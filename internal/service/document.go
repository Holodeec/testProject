package service

import (
	"context"
	"sort"
	"sync"
	"testProject/internal/model"
	"testProject/internal/repository"
	"time"
)

type cachedItem struct {
	doc     *model.Document
	expires time.Time
}

type DocumentService interface {
	Create(ctx context.Context, document *model.Document) error
	FindByID(ctx context.Context, id string) (*model.Document, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.Document, error)
	Update(ctx context.Context, document *model.Document) error
	Delete(ctx context.Context, id string) error
}

type documentService struct {
	documentRepository repository.DocumentRepository
	cache              map[string]cachedItem
	mu                 sync.RWMutex
	ttl                time.Duration
	maxWorkers         int
}

func NewDocumentService(repo repository.DocumentRepository, ttl time.Duration) DocumentService {
	return &documentService{
		documentRepository: repo,
		cache:              make(map[string]cachedItem),
		ttl:                ttl,
		maxWorkers:         10,
	}
}

func (s *documentService) Create(ctx context.Context, document *model.Document) error {
	return s.documentRepository.Create(ctx, document)
}

func (s *documentService) FindAll(ctx context.Context, limit, offset int) ([]*model.Document, error) {
	docs, err := s.documentRepository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	workerCount := s.calculateWorkers(len(docs))

	jobs := make(chan *model.Document, len(docs))
	results := make(chan *model.Document, len(docs))

	s.startWorkers(workerCount, jobs, results)

	for _, doc := range docs {
		jobs <- doc
	}
	close(jobs)

	return s.collectResults(results, len(docs)), nil
}

func (s *documentService) FindByID(ctx context.Context, id string) (*model.Document, error) {
	s.mu.RLock()
	item, found := s.cache[id]
	s.mu.RUnlock()

	if found && time.Now().Before(item.expires) {
		return deepCopyDocument(item.doc), nil
	}

	doc, err := s.documentRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sanitized := sanitizeAndSortCopy(doc)

	s.mu.Lock()
	s.cache[id] = cachedItem{doc: sanitized, expires: time.Now().Add(s.ttl)}
	s.mu.Unlock()

	return deepCopyDocument(sanitized), nil
}

func (s *documentService) Update(ctx context.Context, document *model.Document) error {
	err := s.documentRepository.Update(ctx, document)
	if err != nil {
		return err
	}

	s.mu.Lock()
	delete(s.cache, document.ID)
	s.mu.Unlock()

	return nil
}

func (s *documentService) Delete(ctx context.Context, id string) error {
	err := s.documentRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	s.mu.Lock()
	delete(s.cache, id)
	s.mu.Unlock()

	return nil
}

func sanitizeAndSortCopy(doc *model.Document) *model.Document {
	if doc == nil {
		return nil
	}
	copyDoc := deepCopyDocument(doc)
	sort.SliceStable(copyDoc.Level1, func(i, j int) bool {
		return copyDoc.Level1[i].Sort > copyDoc.Level1[j].Sort
	})
	copyDoc.Description = ""
	for i := range copyDoc.Level1 {
		copyDoc.Level1[i].Name = ""
		for j := range copyDoc.Level1[i].Level2 {
			copyDoc.Level1[i].Level2[j].Value = ""
		}
	}
	return copyDoc
}

func deepCopyDocument(doc *model.Document) *model.Document {
	if doc == nil {
		return nil
	}

	lvl1 := make([]model.Level1Item, len(doc.Level1))
	for i, el1 := range doc.Level1 {
		lvl2 := make([]model.Level2Item, len(el1.Level2))
		for j, el2 := range el1.Level2 {
			lvl2[j] = model.Level2Item{
				Value: el2.Value,
				Code:  el2.Code,
			}
		}
		lvl1[i] = model.Level1Item{
			Name:   el1.Name,
			Sort:   el1.Sort,
			Level2: lvl2,
		}
	}

	return &model.Document{
		ID:          doc.ID,
		Title:       doc.Title,
		Description: doc.Description,
		Level1:      lvl1,
	}
}

func (s *documentService) calculateWorkers(docCount int) int {
	if docCount == 0 {
		return 0
	}
	if s.maxWorkers > docCount {
		return docCount
	}
	return s.maxWorkers
}

func (s *documentService) startWorkers(workerCount int, jobs <-chan *model.Document, results chan<- *model.Document) {
	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			s.worker(jobs, results)
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()
}

func (s *documentService) worker(jobs <-chan *model.Document, results chan<- *model.Document) {
	for doc := range jobs {
		processed := sanitizeAndSortCopy(doc)
		results <- processed
	}
}

func (s *documentService) collectResults(results <-chan *model.Document, expected int) []*model.Document {
	out := make([]*model.Document, 0, expected)
	for res := range results {
		out = append(out, res)
	}
	return out
}

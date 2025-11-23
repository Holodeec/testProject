package repository

import (
	"context"
	"testProject/internal/model"
	app_err "testProject/pkg/app-err"

	"github.com/restream/reindexer/v4"
)

type DocumentRepository interface {
	Create(ctx context.Context, document *model.Document) error
	FindByID(ctx context.Context, id string) (*model.Document, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.Document, error)
	Update(ctx context.Context, document *model.Document) error
	Delete(ctx context.Context, id string) error
}

type documentRepository struct {
	db *reindexer.Reindexer
}

func NewDocumentRepository(db *reindexer.Reindexer) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(ctx context.Context, document *model.Document) error {
	tx, err := r.db.WithContext(ctx).BeginTx("documents")
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.Upsert(document)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *documentRepository) FindByID(ctx context.Context, id string) (*model.Document, error) {
	iterator := r.db.WithContext(ctx).Query("documents").Where("id", reindexer.EQ, id).Exec()
	defer iterator.Close()

	if iterator.Next() {
		return iterator.Object().(*model.Document), nil
	}

	if err := iterator.Error(); err != nil {
		return nil, err
	}

	return nil, app_err.DocumentNotFoundErr
}

func (r *documentRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Document, error) {
	iterator := r.db.WithContext(ctx).Query("documents").Limit(limit).Offset(offset).Exec()
	defer iterator.Close()

	var docs []*model.Document
	for iterator.Next() {
		doc := iterator.Object().(*model.Document)
		copyDoc := *doc
		docs = append(docs, &copyDoc)
	}
	return docs, nil
}

func (r *documentRepository) Update(ctx context.Context, document *model.Document) error {
	tx, err := r.db.WithContext(ctx).BeginTx("documents")
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = r.FindByID(ctx, document.ID)
	if err != nil {
		return err
	}

	err = tx.Update(document)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *documentRepository) Delete(ctx context.Context, id string) error {
	q := r.db.WithContext(ctx).Query("documents").Where("id", reindexer.EQ, id)
	_, err := q.Delete()
	return err
}

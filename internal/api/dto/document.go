package dto

type DocumentRequest struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Level1      []Level1ItemRequest `json:"level1"`
}

type Level1ItemRequest struct {
	Sort   int                 `json:"sort"`
	Name   string              `json:"name"`
	Level2 []Level2ItemRequest `json:"level2"`
}

type Level2ItemRequest struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}

type DocumentResponse struct {
	ID          string               `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Level1      []Level1ItemResponse `json:"level1"`
}

type Level1ItemResponse struct {
	Sort   int                  `json:"sort"`
	Name   string               `json:"name"`
	Level2 []Level2ItemResponse `json:"level2"`
}

type Level2ItemResponse struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}

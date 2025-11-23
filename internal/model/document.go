package model

type Document struct {
	ID          string       `json:"id" reindex:"id,,pk"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Level1      []Level1Item `json:"level1,omitempty"`
}

type Level1Item struct {
	Sort   int          `json:"sort" reindex:"sort,tree"`
	Name   string       `json:"name,omitempty"`
	Level2 []Level2Item `json:"level2,omitempty"`
}

type Level2Item struct {
	Code  string `json:"code,omitempty"`
	Value string `json:"value,omitempty"`
}

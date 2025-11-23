package dto

import (
	"testProject/internal/model"

	"github.com/google/uuid"
)

func ToDocument(req *DocumentRequest) *model.Document {
	return &model.Document{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Level1:      toLvl1(req.Level1),
	}
}

func toLvl1(lvl1 []Level1ItemRequest) []model.Level1Item {
	lvl1Items := make([]model.Level1Item, 0, len(lvl1))
	for _, lvl1Item := range lvl1 {
		l := model.Level1Item{
			Sort:   lvl1Item.Sort,
			Name:   lvl1Item.Name,
			Level2: toLvl2(lvl1Item.Level2),
		}
		lvl1Items = append(lvl1Items, l)
	}
	return lvl1Items
}

func toLvl2(lvl2 []Level2ItemRequest) []model.Level2Item {
	lvl2Items := make([]model.Level2Item, 0, len(lvl2))
	for _, lvl2Item := range lvl2 {
		l := model.Level2Item{
			Code:  lvl2Item.Code,
			Value: lvl2Item.Value,
		}
		lvl2Items = append(lvl2Items, l)
	}
	return lvl2Items
}

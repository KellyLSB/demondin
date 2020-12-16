package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/google/uuid"
	//"github.com/kr/pretty"
)

func CreateItemOption(xo *gorm.DB) *ItemOption {
	var itemOption *ItemOption
	xo.Create(itemOption)
	return itemOption
}

func FetchItemOption(xo *gorm.DB, uuid uuid.UUID) *ItemOption {
	var itemOption ItemOption
	xo.Preload("Items").First(&itemOption, "id = ?", uuid)
	return &itemOption
}

func (o *ItemOption) LoadItemOptionType(xo *gorm.DB) *ItemOption {
	o.ItemOptionType = new(ItemOptionType)
	xo.Model(o).Related(o.ItemOptionType)
	return o
}

func (o *ItemOption) Remove(tx *gorm.DB) {
	tx.Unscoped().Delete(o)
}

func (o *ItemOption) Sample() string {
	return fmt.Sprintf("[%s]: %s",
		o.ItemOptionType.Key,
		o.Values,
	)
}

func (o *ItemOption) Save(tx *gorm.DB) {
	tx.Save(o)
}

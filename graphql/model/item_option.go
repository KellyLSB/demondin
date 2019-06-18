package model

import (
	"fmt"	
	"github.com/jinzhu/gorm"
	"github.com/google/uuid"
	//"github.com/kr/pretty"
)

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

func (o *ItemOption) Sample() string {
	return fmt.Sprintf("[%s]: %s", 
		o.ItemOptionType.Key, 
		o.Values,
	)
}

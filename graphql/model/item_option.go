package model

import (
	"github.com/jinzhu/gorm"
)

func (o *ItemOption) LoadOptionType(tx *gorm.DB) *ItemOption {
	//o.OptionType = new(ItemOptionType)
	tx.Model(o).Related(o.OptionType)
	return o
}

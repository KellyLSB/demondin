package model

import (
	"github.com/jinzhu/gorm"
)

func (i *InvoiceItem) LoadOptions(tx *gorm.DB) *InvoiceItem {
	tx.Model(i).Association("Options").Find(&i.Options)
	return i
}

package model

import (
	//"fmt"
	"github.com/jinzhu/gorm"
	//"github.com/kr/pretty"
)

func (i *InvoiceItem) LoadItem(tx *gorm.DB) *InvoiceItem {
	i.Item = new(Item)	
	tx.Model(i).Related(i.Item)
	return i
}

func (i *InvoiceItem) LoadPrice(tx *gorm.DB) *InvoiceItem {
	i.ItemPrice = new(ItemPrice)
	tx.Model(i).Related(i.ItemPrice)
	return i
}

func (i *InvoiceItem) LoadOptions(tx *gorm.DB) *InvoiceItem {
	tx.Model(i).Related(&i.Options)

	for _, o := range i.Options {
		o.LoadItemOptionType(tx)
	}

	//fmt.Printf("%# v", pretty.Formatter(i))

	return i
}

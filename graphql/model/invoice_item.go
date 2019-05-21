package model

import (
	//"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/google/uuid"
	//"github.com/kr/pretty"
)

func FetchInvoiceItem(tx *gorm.DB, iUUID uuid.UUID) *InvoiceItem {
	var invoiceItem InvoiceItem

	tx.Preload(
		"Item", "ItemPrice", "Options",
	).First(&invoiceItem, "id = ?", iUUID)

	return &invoiceItem	
}

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

func (i *InvoiceItem) AddOption(
	tx *gorm.DB, 
	itemOptionType *ItemOptionType, 
	values postgres.Jsonb,
) (itemOption *ItemOption) {
	itemOption = &ItemOption{
		ItemOptionTypeID: itemOptionType.ID,
		// Get RawJSON from postgres.Jsonb object
		Values: string(values.RawMessage),
	}

	tx.Model(i).Association("Options").Append(itemOption)
	itemOption.LoadItemOptionType(tx)
	
	return
}

func (i *InvoiceItem) AddOptionByTypeUUID(
	tx *gorm.DB,
	itemOptionTypeUUID uuid.UUID,
	values postgres.Jsonb,
) (*ItemOption) {
	return i.AddOption(tx, FetchItemOptionType(
		tx, itemOptionTypeUUID,
	), values)
}

func (i *InvoiceItem) LoadOptions(tx *gorm.DB) *InvoiceItem {
	tx.Model(i).Related(&i.Options)

	for _, o := range i.Options {
		o.LoadItemOptionType(tx)
	}

	//fmt.Printf("%# v", pretty.Formatter(i))

	return i
}

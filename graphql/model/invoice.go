package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/kr/pretty"
)

func FetchInvoice(tx *gorm.DB, inputs ...interface{}) (*Invoice) {
	var invoice *Invoice = new(Invoice)
	var invUUID uuid.UUID

	if len(inputs) < 1 {
		panic("Missing ID Input: FetchInvoice(*gorm.DB, uuid.UUID|string)")
	}

	for _, input := range inputs {
		switch input := input.(type) {
		case Invoice:
			invoice = &input
		case *Invoice:
			invoice = input
		case uuid.UUID:
			invUUID = input
		case string:
			invUUID = uuid.MustParse(input)
		default:
			fmt.Printf("%# v\n", pretty.Formatter(input))
		}
	}
	
	tx.Preload("Items").First(invoice, "id = ?", invUUID)	
	return invoice
}

func FetchOrCreateInvoice(tx *gorm.DB, inputs ...interface{}) (*Invoice) {
	var invoice *Invoice = new(Invoice)
	var invUUID uuid.UUID

	if len(inputs) < 1 {
		goto CreateInvoice
	}

	// Seems like duplication of effort except
	// the variable assignments... think this over
	for _, input := range inputs {
		switch input := input.(type) {
		case Invoice:
			invoice = &input
		case *Invoice:
			invoice = input
		case uuid.UUID:
			invUUID = input
		case string:
			invUUID = uuid.MustParse(input)
		default:
			fmt.Printf("%# v\n", pretty.Formatter(input))
		}
	}

	FetchInvoice(tx, invoice, invUUID)

	if invoice.ID == uuid.Nil {
		goto CreateInvoice
	}
	
	return invoice

CreateInvoice:
	tx.Create(invoice)
	return invoice
}

func (i *Invoice) LoadItems(tx *gorm.DB) *Invoice {
	tx.Model(i).Related(&i.Items)

	for _, it := range i.Items {
		it.LoadItem(tx)
		it.LoadPrice(tx)
		it.LoadOptions(tx)
	}

	//fmt.Printf("%# v", pretty.Formatter(i))

	return i
}

func (i *Invoice) AddItem(tx *gorm.DB, item *Item) (*Item) {
	tx.Model(i).Association("Items").Append(&InvoiceItem{
		InvoiceID: i.ID, ItemID: item.ID,		
		ItemPriceID: item.LoadPrices(tx).CurrentPrice().ID,
	})

	return item
}

func (i *Invoice) AddItemByUUID(tx *gorm.DB, itemUUID uuid.UUID) (*Item) {
	var item Item
	tx.First(&item, "id = ?", itemUUID)
	return i.AddItem(tx, &item)
}

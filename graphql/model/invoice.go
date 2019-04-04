package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func FetchInvoice(tx *gorm.DB, uuid uuid.UUID) (*Invoice) {
	var invoice Invoice
	tx.Preload("Items").First(&invoice, "id = ?", uuid)	
	return &invoice
}

func (i *Invoice) LoadItems(tx *gorm.DB) *Invoice {
	tx.Model(i).Association("Items").Find(&i.Items)
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

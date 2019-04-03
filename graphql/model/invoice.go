package model

import (
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go"
	"github.com/jinzhu/gorm"
)

type Invoice struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   *time.Time     `json:"deletedAt"`
	CardToken   *string        `json:"cardToken"`
	ChargeToken *string        `json:"chargeToken"`
	CardData    *stripe.Card   `json:"cardData"`
	ChargeData  *stripe.Charge `json:"chargeData"`
	Items       []InvoiceItem  `json:"items"`
}

func FetchInvoice(tx *gorm.DB, uuid uuid.UUID) (*Invoice) {
	var invoice Invoice
	tx.Preload("Items").First(&invoice, "id = ?", uuid)	
	return &invoice
}

func (i *Invoice) AddItem(tx *gorm.DB, item *Item) (*Item) {
	fmt.Printf("\n%+v\n", item.CurrentPrice())
	
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

func (Invoice) IsPostgresql() {}

package model

import (
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go"
	"github.com/jinzhu/gorm"
)

type InvoiceItem struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	DeletedAt   *time.Time   `json:"deletedAt"`
	InvoiceID   uuid.UUID    `json:"invoiceID" gorm:"type:uuid"`
	ItemID      uuid.UUID    `json:"itemID" gorm:"type:uuid"`
	ItemPriceID uuid.UUID    `json:"itemPriceID" gorm:"type:uuid"`
	Options     []ItemOption `json:"options"`
}

func (i *InvoiceItem) LoadOptions(tx *gorm.DB) *InvoiceItem {
	tx.Model(i).Association("Options").Find(&i.Options)
	return i
}

func (InvoiceItem) IsPostgresql() {}

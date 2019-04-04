package model

import (
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go"
	"github.com/jinzhu/gorm"
)

type ItemOption struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     *time.Time     `json:"deletedAt"`
	InvoiceItemID uuid.UUID      `json:"invoiceItemID" gorm:"type:uuid"`
	OptionType    ItemOptionType `json:"optionType"`
	Values        string         `json:"values"`
}

func (i *ItemOption) KeyValuePair() {
}

func (ItemOption) IsPostgresql() {}

// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"

	"github.com/KellyLSB/demondin/graphql/postgres"
	"github.com/google/uuid"
)

type Postgresql interface {
	IsPostgresql()
}

type Account struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4();"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	Auth         int        `json:"auth"`
	Name         *string    `json:"name"`
	Email        string     `json:"email" gorm:"unique;"`
	PasswordHash *string    `json:"passwordHash"`
}

func (Account) IsPostgresql() {}

type Invoice struct {
	ID             uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4();"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
	DeletedAt      *time.Time             `json:"deletedAt"`
	AccountID      uuid.UUID              `json:"accountID" gorm:"type:uuid;"`
	Account        *Account               `json:"account"`
	StripeTokenID  *string                `json:"stripeTokenID"`
	StripeChargeID *string                `json:"stripeChargeID"`
	StripeToken    *postgres.StripeToken  `json:"stripeToken" gorm:"type:jsonb;"`
	StripeCharge   *postgres.StripeCharge `json:"stripeCharge" gorm:"type:jsonb;"`
	SubTotal       int                    `json:"subTotal"`
	DemonDin       int                    `json:"demonDin"`
	Taxes          int                    `json:"taxes"`
	Total          int                    `json:"total"`
	Items          []*InvoiceItem         `json:"items"`
}

func (Invoice) IsPostgresql() {}

type InvoiceItem struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4();"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	DeletedAt   *time.Time    `json:"deletedAt"`
	Invoice     *Invoice      `json:"invoice"`
	InvoiceID   uuid.UUID     `json:"invoiceID" gorm:"type:uuid;"`
	Item        *Item         `json:"item" gorm:"save_associations:false;"`
	ItemID      uuid.UUID     `json:"itemID" gorm:"type:uuid;"`
	ItemPrice   *ItemPrice    `json:"itemPrice" gorm:"save_associations:false;"`
	ItemPriceID uuid.UUID     `json:"itemPriceID" gorm:"type:uuid;"`
	Options     []*ItemOption `json:"options"`
}

func (InvoiceItem) IsPostgresql() {}

type ItemOption struct {
	ID               uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4();"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	DeletedAt        *time.Time      `json:"deletedAt"`
	InvoiceItem      *InvoiceItem    `json:"invoiceItem" gorm:"save_associations:false;"`
	InvoiceItemID    uuid.UUID       `json:"invoiceItemID" gorm:"type:uuid;"`
	ItemOptionType   *ItemOptionType `json:"itemOptionType" gorm:"save_associations:false;"`
	ItemOptionTypeID uuid.UUID       `json:"itemOptionTypeID" gorm:"type:uuid;"`
	Values           postgres.Jsonb  `json:"values" gorm:"type:jsonb"`
}

func (ItemOption) IsPostgresql() {}

type ItemOptionType struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4();"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt *time.Time     `json:"deletedAt"`
	Item      *Item          `json:"item" gorm:"save_associations:false;"`
	ItemID    uuid.UUID      `json:"itemID" gorm:"type:uuid;"`
	Key       string         `json:"key"`
	ValueType string         `json:"valueType"`
	Values    postgres.Jsonb `json:"values"`
}

func (ItemOptionType) IsPostgresql() {}

type ItemPrice struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4();"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
	Item       *Item      `json:"item" gorm:"save_associations:false;"`
	ItemID     uuid.UUID  `json:"itemID" gorm:"type:uuid;"`
	Price      int        `json:"price"`
	BeforeDate time.Time  `json:"beforeDate"`
	AfterDate  time.Time  `json:"afterDate"`
	Taxable    bool       `json:"taxable"`
}

func (ItemPrice) IsPostgresql() {}

type NewAccount struct {
	ID       *uuid.UUID `json:"id"`
	Auth     *int       `json:"auth"`
	Name     *string    `json:"name"`
	Email    *string    `json:"email"`
	Password *string    `json:"password"`
}

type NewInvoice struct {
	ID            *uuid.UUID       `json:"id"`
	Account       *NewAccount      `json:"account"`
	StripeTokenID *string          `json:"stripeTokenID"`
	Items         []NewInvoiceItem `json:"items"`
	Submit        *bool            `json:"submit"`
}

type NewInvoiceItem struct {
	ID          *uuid.UUID      `json:"id"`
	ItemID      uuid.UUID       `json:"itemID"`
	ItemPriceID uuid.UUID       `json:"itemPriceID"`
	Options     []NewItemOption `json:"options"`
	Remove      *bool           `json:"remove"`
}

type NewItem struct {
	ID          *uuid.UUID          `json:"id"`
	Name        string              `json:"name"`
	Description *string             `json:"description"`
	Prices      []NewItemPrice      `json:"prices"`
	Options     []NewItemOptionType `json:"options"`
}

type NewItemOption struct {
	ID               *uuid.UUID     `json:"id"`
	ItemOptionTypeID uuid.UUID      `json:"itemOptionTypeID"`
	Values           postgres.Jsonb `json:"values"`
}

type NewItemOptionType struct {
	ID        *uuid.UUID     `json:"id"`
	Key       *string        `json:"key"`
	ValueType *string        `json:"valueType"`
	Values    postgres.Jsonb `json:"values"`
}

type NewItemPrice struct {
	ID         *uuid.UUID `json:"id"`
	Price      int        `json:"price"`
	BeforeDate time.Time  `json:"beforeDate"`
	AfterDate  time.Time  `json:"afterDate"`
}

type Paging struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Session struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4();"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
	AccountID  *uuid.UUID `json:"accountID" gorm:"type:uuid;"`
	Account    *Account   `json:"account" gorm:"save_associations:false;"`
	InvoiceID  *uuid.UUID `json:"invoiceID" gorm:"type:uuid;"`
	Invoice    *Invoice   `json:"invoice" gorm:"save_associations:false;"`
	RemoteAddr *string    `json:"remoteAddr"`
	UserAgent  *string    `json:"userAgent"`
	Referer    *string    `json:"referer"`
	Method     *string    `json:"method"`
	URL        *string    `json:"url"`
}

func (Session) IsPostgresql() {}

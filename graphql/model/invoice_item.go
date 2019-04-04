package model

import (
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go"
	"github.com/jinzhu/gorm"
)

func (i *InvoiceItem) LoadOptions(tx *gorm.DB) *InvoiceItem {
	tx.Model(i).Association("Options").Find(&i.Options)
	return i
}

func (InvoiceItem) IsPostgresql() {}

package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// Item is a general record for items purchasable.
type Item struct {
	Model

	// Enabled for sale
	Enabled bool `gorm:"default:false"`
	IsBadge bool `gorm:"default:false"`

	// Information on Item
	Name        string
	Description []byte
	Options     postgres.Jsonb `gorm:"type:jsonb;"`
}

// Price record should immutable (except for validity)
// as to avoid conflicts in change
type Price struct {
	Model

	// BelongsTo Record
	ItemID uuid.UUID `gorm:"type:uuid;"`

	// Pricing
	Price   int
	Taxable bool

	// When this is valid
	ValidAfter  *time.Time
	ValidBefore *time.Time
}

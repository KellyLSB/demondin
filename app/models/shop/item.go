package shopModel

import (
	"time"

	"github.com/KellyLSB/youmacon/app/models"
	"github.com/google/uuid"
)

// Item is a general record for items purchasable.
type Item struct {
	models.Model

	// Enabled for sale
	Enabled bool

	// Information on Item
	Name        string
	Description []byte
}

// Price record should immutable (except for validity)
// as to avoid conflicts in change
type Price struct {
	models.Model

	// BelongsTo Record
	ItemID uuid.UUID `gorm:"type:uuid;"`

	// Pricing
	Price   int
	Taxable bool

	// When this is valid
	ValidAfter  *time.Time
	ValidBefore *time.Time
}

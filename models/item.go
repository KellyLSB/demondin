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
	Name        string         `gorm:"type:varchar(60);"`
	Description string         `gorm:"type:text;"`
	Options     postgres.Jsonb `gorm:"type:jsonb;"`
	
  // 	{
  // 	  name: string,
  // 	  t-shirt: {
  // 	    type: 'select',
  // 	    multi: [
  // 	      'S', 'M', 'L',
  // 	    ]
  // 	  },
  // 	  otherthing: {
  // 	    type: 'multi-select',
  // 	    multi: [
  // 	      'One', 'Two', 'Three'
  // 	    ]
  // 	  },
  // 	  something: {
  // 	    type: 'radio',
  // 	    multi: [
  // 	      'One', 'Two', 'Three'
  // 	    ]
  // 	  }
  // 	}

	// HasMany
	Prices []Price
}

// Price record should immutable (except for validity)
// as to avoid conflicts in change
type Price struct {
	Model

	// BelongsTo Record
	// perhaps turn off autopreloading for tis item (circular loading xD)
	//Item   Item      `gorm:"foreignkey:item_id;"`
	ItemID uuid.UUID `gorm:"type:uuid;"`

	// Pricing
	Price   int
	Taxable bool

	// When this is valid
	ValidAfter  *time.Time
	ValidBefore *time.Time
}

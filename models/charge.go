package models

import (
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Charge struct {
	Model

	StripeData postgres.Jsonb `gorm:"type:jsonb;"`
}

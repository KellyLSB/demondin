package models

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`

	// Timestamps
	CreatedAt time.Time  `gorm:"index;default:now();"`
	UpdatedAt time.Time  `gorm:"index;default:now();"`
	DeletedAt *time.Time `gorm:"index;default:null;"`
}

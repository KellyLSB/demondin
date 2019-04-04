package model

import (
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go"
	"github.com/jinzhu/gorm"
)
	
func (i *ItemOption) KeyValuePair() {
}

func (ItemOption) IsPostgresql() {}

package postgres

import (
	"database/sql/driver"
	"encoding/json"
	_ "github.com/lib/pq"
	
	"github.com/stripe/stripe-go"
)

type StripeToken struct {
	stripe.Token
}

func (t StripeToken) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *StripeToken) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), t)
}


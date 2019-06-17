package postgres

import (
	"database/sql/driver"
	"encoding/json"
	_ "github.com/lib/pq"
	
	"github.com/stripe/stripe-go"
)

type StripeCharge struct {
	stripe.Charge
}

func (c StripeCharge) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *StripeCharge) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), c)
}


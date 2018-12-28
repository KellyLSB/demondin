package models

import (
  "github.com/google/uuid"
  "encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	stripeClient "github.com/stripe/stripe-go/client"
)

type Invoice struct {
	Model
	
	CardToken   string
	CardData    postgres.Jsonb `gorm:"type:jsonb;"`
	ChargeToken string
  ChargeData  postgres.Jsonb `gorm:"type:jsonb;"`
  
  Badges []*Badge
}

func (i Invoice) GetCardData(stripe *stripeClient.API) {
  card, err := stripe.Cards.Get(i.CardToken, nil)
  if err != nil {
    panic(err)
  }
  
  data, err := json.Marshal(card)
  if err != nil {
    panic(err)
  }
  
  i.CardData = postgres.Jsonb{data}
}

func (i Invoice) GetChargeData(stripe *stripeClient.API) {
  charge, err := stripe.Charges.Get(i.ChargeToken, nil)
  if err != nil {
    panic(err)
  }
  
  data, err := json.Marshal(charge)
  if err != nil {
    panic(err)
  }
  
  i.ChargeData = postgres.Jsonb{data}
}

type Badge struct {
  Model
  
  // Item that badge is sourced from.
  Item   Item      `gorm:"foreignkey:ItemID;"`
  ItemID uuid.UUID `gorm:"type:uuid;"`
  
  // Price item is bought at.
  Price   Price     `gorm:"foreignkey:PriceID;"`
  PriceID uuid.UUID `gorm:"type:uuid;"`
  
  // Disabled due to preload looping
  // Invoice Invoice `gorm:"foreign_key:invoice_id;"`
  InvoiceID uuid.UUID `gorm:"type:uuid;"`
  
  Options postgres.Jsonb `gorm:"type:jsonb;"`
  // {
  //   name: "KellyLSBN",
  //   t-shirt: 'S',
  //   otherthing: ['One', 'Two'],
  //   something: 'One'
  // }
}
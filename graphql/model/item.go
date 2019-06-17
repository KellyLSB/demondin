package model

import (
	//"fmt"
	//"log"
	//"sort"
	"sort"
	"time"
	"io"
	"strconv"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/araddon/dateparse"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/99designs/gqlgen/graphql"
	"github.com/jinzhu/gorm"

	"github.com/stripe/stripe-go"
)

type Item struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
	DeletedAt   *time.Time       `json:"deletedAt"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	Enabled     bool             `json:"enabled"`
	IsBadge     bool             `json:"isBadge"`
	Options     []ItemOptionType `json:"options" gorm:"foreignkey:ItemID"`
	Prices      []ItemPrice      `json:"prices" gorm:"foreignkey:ItemID"`
}

func FetchItem(tx *gorm.DB, uuid uuid.UUID) *Item {
	var item Item	
	tx.First(&item, "id = ?", uuid)
	return &item
}

func (i *Item) LoadOptions(tx *gorm.DB) *Item {
	tx.Model(i).Association("Options").Find(&i.Options)
	return i
}

func (i *Item) LoadPrices(tx *gorm.DB) *Item {
	tx.Model(i).Association("Prices").Find(&i.Prices)
	return i
}

func (i *Item) CurrentPrice() *ItemPrice {
	sort.Slice(i.Prices, func(x, y int) bool {
		sb :=        i.Prices[x].AfterDate.After(i.Prices[y].AfterDate)
		return sb && i.Prices[x].BeforeDate.Before(i.Prices[y].BeforeDate)
	})

	now := time.Now()
	for _, price := range i.Prices {
		if now.Before(price.AfterDate) || now.After(price.BeforeDate) {
			continue
		}
		
		return &price
	}

	return nil
}

func (Item) IsPostgresql() {}

func MarshalID(id uuid.UUID) graphql.Marshaler {
  return graphql.WriterFunc(func(w io.Writer) {
    io.WriteString(w, strconv.Quote(id.String()))
  })
}

func UnmarshalID(v interface{}) (uuid.UUID, error) {
  return uuid.Parse(v.(string))
}

func MarshalDateTime(t time.Time) graphql.Marshaler {
  return graphql.WriterFunc(func(w io.Writer) {
    io.WriteString(w, strconv.Quote(t.Format(time.RFC3339)))
  })
}

func UnmarshalDateTime(v interface{}) (time.Time, error) {
  return dateparse.ParseAny(v.(string))
}

func MarshalJSON(v postgres.Jsonb) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		byt, _ := v.MarshalJSON()
		w.Write(byt)
	})
}

func UnmarshalJSON(v interface{}) (out postgres.Jsonb, err error) {
	byt, err := json.Marshal(v)
	err = out.UnmarshalJSON(byt)
	return
}

// @TODO:
// Might break this into specific scalars that serialize the
// official models: https://gowalker.org/github.com/stripe/stripe-go#Card
func MarshalStripeToken(v stripe.Token) graphql.Marshaler {
  return graphql.WriterFunc(func(w io.Writer) {
    json.NewEncoder(w).Encode(&v)
  })
}

func UnmarshalStripeToken(v interface{}) (out stripe.Token, err error) {
  err = json.NewDecoder(strings.NewReader(v.(string))).Decode(&out)
  return
}

func MarshalStripeCharge(v stripe.Charge) graphql.Marshaler {
  return graphql.WriterFunc(func(w io.Writer) {
    json.NewEncoder(w).Encode(&v)
  })
}

func UnmarshalStripeCharge(v interface{}) (out stripe.Charge, err error) {
  err = json.NewDecoder(strings.NewReader(v.(string))).Decode(&out)
  return
}


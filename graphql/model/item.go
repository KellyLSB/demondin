package model

import (
  //"log"
  //"sort"
	"time"
	//"fmt"
	"io"
	"strconv"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/araddon/dateparse"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/99designs/gqlgen/graphql"

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
    io.WriteString(w, string(v.RawMessage))
  })
}

func UnmarshalJSON(v interface{}) (out postgres.Jsonb, err error) {
  byt, err := json.Marshal(v)
  return postgres.Jsonb{json.RawMessage(byt)}, err
}

// @TODO:
// Might break this into specific scalars that serialize the
// official models: https://gowalker.org/github.com/stripe/stripe-go#Card
func MarshalStripeCard(v stripe.Card) graphql.Marshaler {
  return graphql.WriterFunc(func(w io.Writer) {
    json.NewEncoder(w).Encode(v)
  })
}

func UnmarshalStripeCard(v interface{}) (out stripe.Card, err error) {
  err = json.NewDecoder(strings.NewReader(v.(string))).Decode(&out)
  return
}

func MarshalStripeCharge(v stripe.Charge) graphql.Marshaler {
  return graphql.WriterFunc(func(w io.Writer) {
    json.NewEncoder(w).Encode(v)
  })
}

func UnmarshalStripeCharge(v interface{}) (out stripe.Charge, err error) {
  err = json.NewDecoder(strings.NewReader(v.(string))).Decode(&out)
  return
}


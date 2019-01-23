package model

import (
  //"log"
  //"sort"
	"time"
	//"fmt"
	"io"
	"strconv"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/araddon/dateparse"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/99designs/gqlgen/graphql"
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
	Options     []ItemOptionType `json:"options"`
	Prices      []ItemPrice      `json:"prices"`
}

func (Item) IsPostgresql() {}

func MarshalID(id uuid.UUID) graphql.Marshaler {
  return graphql.WriterFunc(func(w io.Writer) {
    io.WriteString(w, id.String())
  })
}

func UnmarshalID(v interface{}) (uuid.UUID, error) {
  str, ok := v.(string)
  if !ok {
    panic("ID must be a string")
  }
  return uuid.Parse(str)
}

func MarshalDateTime(t time.Time) graphql.Marshaler {
  return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatInt(t.Unix(), 10))
	})
}

func UnmarshalDateTime(v interface{}) (time.Time, error) {
  return dateparse.ParseAny(v.(string))
}

func MarshalJSON(v postgres.Jsonb) graphql.Marshaler {
  return graphql.WriterFunc(func(w io.Writer) {
    json.NewEncoder(w).Encode(v)
  })
}

func UnmarshalJSON(v interface{}) (out postgres.Jsonb, err error) {
  byt, err := json.Marshal(v)
  if err != nil {
    return
  }
  
  return postgres.Jsonb{json.RawMessage(byt)}, nil
}
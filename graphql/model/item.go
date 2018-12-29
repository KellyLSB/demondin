package model

import (
  //"log"
  //"sort"
	"time"
	"fmt"
	"io"
	"strconv"

	"github.com/google/uuid"
	//"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/99designs/gqlgen/graphql"
)

type Item struct {
	ID          uuid.UUID        `json:"id"`
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
  if tmpStr, ok := v.(int64); ok {
		return time.Unix(tmpStr, 0), nil
	}
	return time.Time{}, fmt.Errorf("time should be a unix timestamp")
}
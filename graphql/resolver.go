package graphql

import (
  "io"
  "os"
  "fmt"
	"context"
	"encoding/json"

	"github.com/KellyLSB/demondin/graphql/model"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var dbh func(func(*gorm.DB))

func init() {
  // Load .env file from repo
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	
  dbh = dbInit(
    os.Getenv("POSTGRES_HOSTPORT"), os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"),
	)
}

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateItem(
  ctx context.Context,
  input model.NewItem,
) (
  item model.Item,
  err error,
) {
  err = pipeInput(&input, &item)
  
  if err != nil {
    return
  }
  
  dbh(func(db *gorm.DB) {
    query := db.Create(&item)
    err = gormErrors(query)
  })
  
  return
}

func (r *mutationResolver) UpdateItem(
  ctx context.Context,
  id uuid.UUID,
  input model.NewItem,
) (
  item model.Item,
  err error,
) {
	dbh(func(db *gorm.DB) {
	  query := db.First(&item, id)
	  err = gormErrors(query)
	  fmt.Printf("%s\n", err)
	})
	
	if err != nil {
	  return
	}
	
	err = pipeInput(&input, &item)
	
	if err != nil {
	  return
	}
	
	dbh(func(db *gorm.DB) {
	  query := db.Save(&item)
	  err = gormErrors(query)
	})
	
	return
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Items(ctx context.Context, paging *model.Paging) ([]model.Item, error) {
  var models []model.Item
  var err    error
  
	dbh(func(db *gorm.DB) {
	  query := db.Select("*").Table("items")
	    
    if paging.Limit > 0 {
	    query = query.Limit(paging.Limit)
    }

    if paging.Offset > 0 {
      query = query.Offset(paging.Offset)
    }
    
    // Add the owner ID's in ItemPrices and ItemOptionTypes
    query = query.Preload("Prices").Preload("Options")
	    
	  query.Find(&models)
	  err = gormErrors(query)
	})
	  
	return models, err
}

func (r *queryResolver) Invoices(ctx context.Context, paging *model.Paging) ([]model.Invoice, error) {
	panic("not implemented")
}

func gormErrors(db *gorm.DB) (err error) {
  for _, e := range db.GetErrors() {
	  err = fmt.Errorf("%s\n%s", err, e)
	}
	
	return
}

func pipeInput(from, to interface{}) (err error) {
  jReader, jWriter := io.Pipe()
  defer jReader.Close()
  
  fmt.Printf(">> Trace JSON Pipe\n")
  defer fmt.Printf("<< Trace JSON Pipe\n")
  
  go func() {
    defer jWriter.Close()
	  err = json.NewEncoder(jWriter).Encode(from)
	}()
	
	if err != nil {
	  return
	}
	
	return json.NewDecoder(jReader).Decode(to)
}
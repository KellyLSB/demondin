package graphql

import (
  "io"
  "os"
  "fmt"
  "time"
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
  // Copy input into the item
  err = pipeInput(&input, &item)
  
  if err != nil {
    return
  }
  
  // Save the record in DB
  dbh(func(db *gorm.DB) {
    err = gormErrors(db.Create(&item))
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
  // Fetch first item by UUID
	dbh(func(db *gorm.DB) {
	  err = gormErrors(db.First(&item, "id = ?", id))
	})
	
	if err != nil {
	  return
	}
	
	// Copy input into the item
	// @TODO: verify updating with\/out associations
	err = pipeInput(&input, &item)
	
	if err != nil {
	  return
	}
	
	// Save the resulting model
	dbh(func(db *gorm.DB) {
	  err = gormErrors(db.Save(&item))
	})
	
	return
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Items(ctx context.Context, paging *model.Paging) ([]model.Item, error) {
  var models []model.Item
  var err    error
  
	dbh(func(db *gorm.DB) {
		query := gormPaging(db.Select("*").Table("items"), paging)
		query = query.Preload("Options")
    
   		// Admin Check
		if false {
			query = query.Preload("Prices")
		} else {
			query = query.Preload("Prices", 
				"item_prices.price > 0 AND " +
				 "? BETWEEN item_prices.after_date" + 
				 " AND item_prices.before_date", time.Now())
		}
                
		err = gormErrors(query.Find(&models))
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

func gormPaging(query *gorm.DB, paging *model.Paging) (*gorm.DB) {
  if paging.Limit > 0 {
    query = query.Limit(paging.Limit)
  }

  if paging.Offset > 0 {
    query = query.Offset(paging.Offset)
  }
  
  return query
}

func pipeInput(from, to interface{}) (err error) {
  jReader, jWriter := io.Pipe()
  defer jReader.Close()
  
  go func() {
    defer jWriter.Close()
	  err = json.NewEncoder(jWriter).Encode(from)
	}()
	
	if err != nil {
	  return
	}
	
	return json.NewDecoder(jReader).Decode(to)
}

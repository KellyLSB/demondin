package graphql

import (
  "os"
  "fmt"
	"context"

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

func (r *mutationResolver) CreateItem(ctx context.Context, input model.NewItem) (model.Item, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateItem(ctx context.Context, id uuid.UUID, input model.NewItem) (model.Item, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Items(ctx context.Context, paging *model.Paging) ([]model.Item, error) {
  var models []model.Item
  var err    error
  
	dbh(func(db *gorm.DB) {
	  query := db.Select("*").
	    Table("items")
	    
    if paging.Limit > 0 {
	    query = query.Limit(paging.Limit)
    }

    if paging.Offset > 0 {
      query = query.Offset(paging.Offset)
    }
	    
	  query.Find(&models)
	    
	  fmt.Printf("%+v\n", models)
	    
	  for _, e := range query.GetErrors() {
	    err = fmt.Errorf("%s\n%s", err, e)
	  }
	  
	  fmt.Printf("%s\n", err)
	})
	  
	return models, err
}

func (r *queryResolver) Invoices(ctx context.Context, paging *model.Paging) ([]model.Invoice, error) {
	panic("not implemented")
}

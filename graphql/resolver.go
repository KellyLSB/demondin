package graphql

import (
  "io"
  "os"
  "fmt"
  "time"
	"context"
	"encoding/json"

	//"github.com/99designs/gqlgen/graphql"
	"github.com/KellyLSB/demondin/graphql/model"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/go-macaron/session"
)

var dbh func(func(*gorm.DB))

var _db Database

func init() {
  	// Load .env file from repo
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	
	_db = InitDB(
		os.Getenv("POSTGRES_HOSTPORT"), os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"),
	)

  	dbh = func(fn func(db *gorm.DB)) {
		_db.Migrate()
		_db.Transact(fn)
	}
}

type subscriptions struct {
	Invoice []chan *model.Invoice
}

var Subscriptions = new(subscriptions)


// Is the resolver a global constant when initiatied by NewExecutableSchema
// or will these subscription bindings be for loss placed here.
// Ideally I will use PG Notify...
type Resolver struct{
	Session session.Store
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

//#  ###   #
//#    #  #
//###  #    #


//#
//# Create / Update process is near identical at the core
//# should consider some wrapp(er|ing) function to streamline process.


type mutationResolver struct{ *Resolver }


func (r *mutationResolver) CreateItem(
	ctx context.Context,
	input model.NewItem,
) (
	*model.Item,
	error,
) {
 	var item model.Item
	var err error

  	// Copy input into the item
	err = pipeInput(&input, &item)	
	if err != nil {
    		return nil, err
 	}
  
	// Save the record in DB
	dbh(func(db *gorm.DB) {
		err = gormErrors(db.Create(&item))
	})
  
	return &item, err
}


func (r *mutationResolver) UpdateItem(
	ctx context.Context,
	id uuid.UUID,
	input model.NewItem,
) (
  	*model.Item,
  	error,
) {
	var item model.Item
	var err error

	// Fetch first item by UUID
	dbh(func(db *gorm.DB) {
		err = gormErrors(db.First(&item, "id = ?", id))
	})
	
	if err != nil {
		return nil, err
	}
	
	// Copy input into the item
	// @TODO: verify updating with\/out associations
	err = pipeInput(&input, &item)
	
	if err != nil {
		return nil, err
	}
	
	// Save the resulting model
	dbh(func(db *gorm.DB) {
		err = gormErrors(db.Save(&item))
	})
	
	return &item, err
}


func (r *mutationResolver) CreateInvoice(
	ctx context.Context, 
	input model.NewInvoice,
) (
	*model.Invoice, 
	error,
) {
	var invoice model.Invoice
	var err error

	// Copy input into the invoice
	err = pipeInput(&input, &invoice)
	if err != nil {
		return nil, err
	}
  
	// Save the record in DB
	dbh(func(db *gorm.DB) {
		err = gormErrors(db.Create(&invoice))
	})

	// Inform subscriptions of create
	fmt.Printf("%d Invoice Subscriptions\n", len(Subscriptions.Invoice))
	for _, sub := range Subscriptions.Invoice {
		sub <- &invoice
	}
  
	return &invoice, err
}


func (r *mutationResolver) UpdateInvoice(
	ctx context.Context, 
	id uuid.UUID, 
	input model.NewInvoice,
) (
	invoice *model.Invoice, 
	err error,
) {
	// Fetch first invoice by UUID
	dbh(func(db *gorm.DB) {
	  err = gormErrors(db.First(invoice, "id = ?", id))
	})
	
	if err != nil {
	  return
	}
	
	// Copy input into the invoice
	// @TODO: verify updating with\/out associations
	err = pipeInput(&input, invoice)
	
	if err != nil {
	  return
	}
	
	// Save the resulting model
	dbh(func(db *gorm.DB) {
	  err = gormErrors(db.Save(invoice))
	})

	// Inform subscriptions of update
	fmt.Printf("%d Invoice Subscriptions\n", len(Subscriptions.Invoice))
	for _, sub := range Subscriptions.Invoice {
		sub <- invoice
	}
	
	return
}

func (r *mutationResolver) AddItemToInvoice(
	ctx context.Context, 
	invoiceID, itemID uuid.UUID, 
	input postgres.Jsonb,
) (
	invoice *model.Invoice, 
	err error,
) {
	var item 	model.Item	
	
	dbh(func(db *gorm.DB) {
		err = gormErrors(db.First(&item, "id = ?", itemID))
	})
	
	if err != nil {
		return
	}

	//invoiceItem := model.InvoiceItem{
	//	InvoiceID: invoiceID, ItemID: itemID,		
	//	ItemPriceID: item.CurrentPrice().ID,		
	//}
	

	// Fetch first invoice by UUID
	dbh(func(db *gorm.DB) {
	  err = gormErrors(db.First(invoice, "id = ?", invoiceID))
	})
	
	if err != nil {
	  return
	}
	
	// Copy input into the invoice
	// @TODO: verify updating with\/out associations
	err = pipeInput(&input, invoice)
	
	if err != nil {
	  return
	}
	
	// Save the resulting model
	dbh(func(db *gorm.DB) {
	  err = gormErrors(db.Save(invoice))
	})

	// Inform subscriptions of update
	fmt.Printf("%d Invoice Subscriptions\n", len(Subscriptions.Invoice))
	for _, sub := range Subscriptions.Invoice {
		sub <- invoice
	}
	
	return
}

//# #    #  ##
//##        ##
//#  #     # #

type queryResolver struct{ *Resolver }


func (r *queryResolver) Items(
	ctx context.Context, 
	paging *model.Paging,
) (
	items []model.Item, 
	err error,
) {
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
                
		err = gormErrors(query.Find(&items))
	})
	  
	return
}


func (r *queryResolver) Invoices(
	ctx context.Context, 
	paging *model.Paging,
) (
	invoices []model.Invoice, 
	err error,
) {
	panic("not implemented")
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) InvoiceUpdated(
	ctx context.Context, 
	id uuid.UUID,
) (<-chan *model.Invoice, error) {
	ch := make(chan *model.Invoice)
	Subscriptions.Invoice = append(Subscriptions.Invoice, ch)
	fmt.Printf("%d Subscriptions Created\n", len(Subscriptions.Invoice))

	fmt.Println("Creating Channel for InvoiceUpdated Subscription.")
	go func() {
		<-ctx.Done()
		fmt.Println("Context for InvoiceUpdated Subscription is Done?!")
	}()
	
	// https://www.postgresql.org/docs/current/sql-createpublication.html
	// https://www.apollographql.com/docs/react/advanced/subscriptions.html
	// https://medium.freecodecamp.org/https-medium-com-anshap1719-graphql-subscriptions-with-go-gqlgen-and-mongodb-5e008fc46451
	// https://graphql.org/blog/subscriptions-in-graphql-and-relay/

	return ch, nil
}

//#   # #  #
//###  #
//#  # #

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

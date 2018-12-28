package main

import (
  "log"
  "time"
  
  "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gopkg.in/macaron.v1"
  "github.com/go-macaron/session"

	"github.com/KellyLSB/demondin/models"
	//stripeClient "github.com/stripe/stripe-go/client"
)

type Invoicing struct {
  Invoice models.Invoice
}

func (i *Invoicing) LoadSession(
  sess session.Store,
  db *gorm.DB, log *log.Logger,
) {
  invoiceUUID := sess.Get("invoice_uuid")
  if invoiceUUID != nil && invoiceUUID != "" {
    log.Printf("Loading Invoice: %s\n", invoiceUUID)
    db.Table("invoices").Preload("Badges").Select("*").
	    Where("invoices.id = ?", invoiceUUID).First(&i.Invoice)
	  return
  }
  
  log.Printf("Initializing a new Invoice\n")
}

func (i *Invoicing) PostAddToCart(
  ctx *macaron.Context,
  db *gorm.DB, log *log.Logger,
) {
  var itemUUID uuid.UUID
  var item models.Item
  JSONDecoder(ctx, &itemUUID)
  
  log.Printf("Loading Item: %s\n", itemUUID)
  db.Table("items").Select("*").Where("items.is_badge AND items.enabled").
    Preload("Prices", "? BETWEEN prices.valid_after AND prices.valid_before", time.Now()).
    Where("items.id = ?", itemUUID).First(&item)
  
  i.Invoice.Badges = append(
    i.Invoice.Badges, &models.Badge{
      Item: item,
      Price: item.CurrentPrice(),
    },
  )
  
  log.Printf("%+v\n", i.Invoice)
    
  log.Printf("ItemID: %s\n", itemUUID)
}

func (i *Invoicing) SaveSession(
  sess session.Store,
  db *gorm.DB, log *log.Logger,
) {
  db.Save(&i.Invoice)
  sess.Set("invoice_uuid", i.Invoice.ID)
  log.Printf("Saving Invoice: %s\n", i.Invoice.ID)
}
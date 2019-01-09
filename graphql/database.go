package graphql

import (
	"fmt"
	"net"

	"github.com/KellyLSB/demondin/graphql/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm"
)

func dbInit(
	hostport, database,
	username, password string,
) func(func(*gorm.DB)) {
  return func(fn func(*gorm.DB)) {
  	// Split the hostport
  	host, port, err := net.SplitHostPort(hostport)
  	if err != nil {
  		panic(err)
  	}
  
  	// Connect to the database
  	db, err := gorm.Open("postgres", fmt.Sprintf(
  		"host=%s port=%s user=%s dbname=%s password=%s",
  		host, port, username, database, password,
  	))
  	if err != nil {
  		panic(fmt.Errorf("DB Connection Failed: %s", err))
  	}
  
  	// Shutdown DB after request
  	defer db.Close()
  
  	// Set logger and migrate
  	db.LogMode(true)
  	dbMigrate(db)
  	fn(db)
  }
}

func dbMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.Invoice{})
	db.AutoMigrate(&model.InvoiceItem{})
	db.AutoMigrate(&model.Item{})
	db.AutoMigrate(&model.ItemOption{})
	db.AutoMigrate(&model.ItemOptionType{})
	db.AutoMigrate(&model.ItemPrice{})
}

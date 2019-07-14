package graphql

import (
	"fmt"
	"net"
	"time"

	"github.com/KellyLSB/demondin/graphql/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Database struct {
	Host, Port, Database,
	Username, Password string

	*gorm.DB
}

func (d Database) Listen() *pq.Listener {
	return pq.NewListener(
		d.ConnInfo(), 10*time.Second, time.Minute, 
		func(ev pq.ListenerEventType, err error) {
			if err != nil {
				fmt.Println(err.Error())
			}
	}	)
}

func (d Database) Transact(fn func(*gorm.DB)) {
	db := d.Open()
	defer db.Close()

	fn(db)
}

func (d Database) Open() (*gorm.DB) {
	var err error

	d.DB, err = gorm.Open("postgres", d.ConnInfo())
	defer d.DB.LogMode(true)
	
	if err != nil {
		panic(err)
	}

	return d.DB
}

func (d Database) ConnInfo() string {
	return fmt.Sprintf(
  		"host=%s port=%s user=%s dbname=%s password=%s",
  		d.Host, d.Port, d.Username, d.Database, d.Password,
  	)
}

func InitDB(
	hostport, database, 
	username, password string,
) Database {
  	host, port, err := net.SplitHostPort(hostport)
  	if err != nil {
  		panic(err)
  	}

	return Database {
		Host: host, Port: port, Database: database,
		Username: username, Password: password,
	}
}

func (d Database) Migrate() {
	d.Transact(func(db *gorm.DB) {
		db.AutoMigrate(&model.Session{})
		db.AutoMigrate(&model.Account{})
		db.AutoMigrate(&model.Invoice{})
		db.AutoMigrate(&model.InvoiceItem{})
		db.AutoMigrate(&model.Item{})
		db.AutoMigrate(&model.ItemOption{})
		db.AutoMigrate(&model.ItemOptionType{})
		db.AutoMigrate(&model.ItemPrice{})
	})
}

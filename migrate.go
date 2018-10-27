package main

import (
	"fmt"
	"log"
	"net"

	"github.com/KellyLSB/demondin/models"
	"github.com/jinzhu/gorm"
	macaron "gopkg.in/macaron.v1"
)

func dbInit(
	hostport, database,
	username, password string,
) macaron.Handler {
	return func(ctx *macaron.Context, logger *log.Logger) {
		// Split the hostport
		host, port, err := net.SplitHostPort(hostport)
		if err != nil {
			logger.Fatal(err)
		}

		// Connect to the database
		db, err := gorm.Open("postgres", fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s",
			host, port, username, database, password,
		))
		if err != nil {
			logger.Fatalf("DB Connection Failed: %s", err)
		}

		// Shutdown DB after request
		defer db.Close()

		// Set logger and migrate
		db.LogMode(true)
		db.SetLogger(logger)
		dbMigrate(db)

		// Map the DB Instance
		ctx.Map(db)
		ctx.Next()
	}
}

func dbMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.M2M{})
	db.AutoMigrate(&models.File{})
	db.AutoMigrate(&models.Item{})
	db.AutoMigrate(&models.Price{})
	db.AutoMigrate(&models.Charge{})
}

package controllers

import (
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Controller is the inherited http controller structure.
type Controller struct {
	DB *gorm.DB
}

// Panic should eventually implement render logic
// that prints to either or both the console or http.Response.
func (c Controller) Panic(err error) {
	if err != nil {
		panic(err)
	}
}

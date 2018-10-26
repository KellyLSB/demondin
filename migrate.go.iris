package main

import (
	"github.com/KellyLSB/youmacon/app/models"
	"github.com/jinzhu/gorm"

	"github.com/KellyLSB/youmacon/app/models/shop"
)

func dbMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.M2M{})
	db.AutoMigrate(&models.File{})
	db.AutoMigrate(&shopModel.Item{})
	db.AutoMigrate(&shopModel.Price{})
	db.AutoMigrate(&shopModel.Item{})
}

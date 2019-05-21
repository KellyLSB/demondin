package model

import (
	//"fmt"	
	"github.com/jinzhu/gorm"
	"github.com/google/uuid"
	//"github.com/kr/pretty"
)

func FetchItemOptionType(xo *gorm.DB, uuid uuid.UUID) *ItemOptionType {
	var itemOptionType ItemOptionType
	xo.First(&itemOptionType, "id = ?", uuid)	
	return &itemOptionType
}

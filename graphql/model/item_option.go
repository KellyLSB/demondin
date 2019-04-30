package model

import (
	//"fmt"	
	"github.com/jinzhu/gorm"
	//"github.com/kr/pretty"
)

func (o *ItemOption) LoadItemOptionType(xo *gorm.DB) *ItemOption {
	o.ItemOptionType = new(ItemOptionType)
	xo.Model(o).Related(o.ItemOptionType)
	//fmt.Printf("% #v", pretty.Formatter(o))
	return o
}

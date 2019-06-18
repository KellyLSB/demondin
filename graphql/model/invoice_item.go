package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/KellyLSB/demondin/graphql/postgres"
	"github.com/google/uuid"
	//"github.com/kr/pretty"
)

func FetchInvoiceItem(tx *gorm.DB, iUUID uuid.UUID) *InvoiceItem {
	var invoiceItem InvoiceItem

	tx.Preload(
		"Item", "ItemPrice", "Options",
	).First(&invoiceItem, "id = ?", iUUID)

	return &invoiceItem
}

func (i *InvoiceItem) LoadItem(tx *gorm.DB) *InvoiceItem {
	i.Item = new(Item)
	tx.Model(i).Related(i.Item)
	return i
}

func (i *InvoiceItem) LoadPrice(tx *gorm.DB) *InvoiceItem {
	i.ItemPrice = new(ItemPrice)
	tx.Model(i).Related(i.ItemPrice)
	return i
}

func (i *InvoiceItem) LoadOptions(tx *gorm.DB) *InvoiceItem {
	tx.Model(i).Related(&i.Options)

	for _, o := range i.Options {
		o.LoadItemOptionType(tx)
	}

	return i
}

func (i *InvoiceItem) LoadRelations(tx *gorm.DB) *InvoiceItem {
	i.LoadOptions(tx)
	i.LoadPrice(tx)
	i.LoadItem(tx)
	return i
}

func (i *InvoiceItem) AddItemOption(
	tx *gorm.DB, 
	itemOption *ItemOption,
) (*ItemOption) {
	tx.Model(i).Association("Options").Append(itemOption)
	return itemOption
}

func (i *InvoiceItem) AddItemOptionType(
	tx *gorm.DB, 
	itemOptionType *ItemOptionType, 
	values postgres.Jsonb,
) (itemOption *ItemOption) {
	return i.AddItemOption(tx, &ItemOption{
		ItemOptionTypeID: itemOptionType.ID,
		// Get RawJSON from postgres.Jsonb object
		Values: postgres.Jsonb{ values.RawMessage },
	})
}

func (i *InvoiceItem) AddItemOptionTypeByUUID(
	tx *gorm.DB,
	itemOptionTypeUUID uuid.UUID,
	values postgres.Jsonb,
) (*ItemOption) {
	return i.AddItemOptionType(tx, FetchItemOptionType(
		tx, itemOptionTypeUUID,
	), values)
}


func (i *InvoiceItem) Sample() string {
	return fmt.Sprintf("%q (%f)", 
		i.Item.Name,
		float32(i.ItemPrice.Price) / 100,
	)
}

func (i *InvoiceItem) SampleOptions() []string {
		var options = []string{}

		for _, _option_ := range i.Options {
			options = append(options, 
				_option_.Sample(),
			)
		} 

		return options
}

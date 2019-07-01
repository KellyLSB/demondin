package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/KellyLSB/demondin/graphql/postgres"
	"github.com/KellyLSB/demondin/graphql/utils"
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

func (i *InvoiceItem) AddItemOption(tx *gorm.DB, io *ItemOption) *ItemOption {
	tx.Model(i).Association("Options").Append(io)
	return io
}

func (i *InvoiceItem) Input(tx *gorm.DB, input *NewInvoiceItem) {
	for _, option := range input.Options {
		i.AddItemOption(tx, &ItemOption{
			ID: utils.EnsureUUID(option.ID),
			ItemOptionTypeID: option.ItemOptionTypeID,
			// Get RawJSON from postgres.Jsonb object
			Values: postgres.Jsonb{ option.Values.RawMessage },
		})
	}
}

// Remove performs an Unsafe Delete, assuming cart is unchecked out.
func (i *InvoiceItem) Remove(tx *gorm.DB) {
	i.LoadOptions(tx)
	
	for _, option := range i.Options {
		option.Remove(tx)
	}
	
	tx.Unscoped().Delete(i)
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

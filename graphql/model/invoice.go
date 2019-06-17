package model

import (
	"os"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/kr/pretty"
	"github.com/KellyLSB/demondin/graphql/utils"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/token"
	"github.com/stripe/stripe-go/charge"
	"gopkg.in/yaml.v2"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_PRIVATE_KEY")
}

func FetchInvoice(tx *gorm.DB, inputs ...interface{}) *Invoice {
	var invoice *Invoice = new(Invoice)
	var invUUID uuid.UUID

	if len(inputs) < 1 {
		panic("Missing ID Input: FetchInvoice(*gorm.DB, uuid.UUID|string)")
	}

	for _, input := range inputs {
		switch input := input.(type) {
		case Invoice:
			invoice = &input
		case *Invoice:
			invoice = input
		case uuid.UUID:
			if input != uuid.Nil {
				invUUID = input
			}
		case *uuid.UUID:
			if input != nil {
				invUUID = *input
			}
		case string:
			invUUID = uuid.MustParse(input)
		default:
			fmt.Printf("%# v\n", pretty.Formatter(input))
		}
	}
	
	tx.Preload("Items").First(invoice, "id = ?", invUUID)	
	return invoice
}

func FetchOrCreateInvoice(tx *gorm.DB, inputs ...interface{}) *Invoice {
	var invoice *Invoice = new(Invoice)
	var invUUID uuid.UUID

	if len(inputs) < 1 {
		goto CreateInvoice
	}

	// Seems like duplication of effort except
	// the variable assignments... think this over
	for _, input := range inputs {
		switch input := input.(type) {
		case Invoice:
			invoice = &input
		case *Invoice:
			invoice = input
		case uuid.UUID:
			if input != uuid.Nil {
	fmt.Printf("%# v\n", pretty.Formatter(input))
				invUUID = input
			}
		case *uuid.UUID:
			if input != nil {
				invUUID = *input
			}
		case string:
			invUUID = uuid.MustParse(input)
		default:
			fmt.Printf("%# v\n", pretty.Formatter(input))
		}
	}

	fmt.Printf("%# v\n", pretty.Formatter(invUUID))

	FetchInvoice(tx, invoice, invUUID)

	if invoice.ID == uuid.Nil {
		goto CreateInvoice
	}
	
	return invoice

CreateInvoice:
	tx.Create(invoice)
	return invoice
}

func (i *Invoice) LoadItems(tx *gorm.DB) *Invoice {
	tx.Model(i).Related(&i.Items)

	for _, it := range i.Items {
		it.LoadItem(tx)
		it.LoadPrice(tx)
		it.LoadOptions(tx)
	}

	//fmt.Printf("%# v", pretty.Formatter(i))

	return i
}

func (i *Invoice) AddInvoiceItem(
	tx *gorm.DB, invoiceItem *InvoiceItem,
) *InvoiceItem {
	tx.Model(i).Association("Items").Append(invoiceItem)
	invoiceItem.LoadOptions(tx)
	invoiceItem.LoadPrice(tx)
	invoiceItem.LoadItem(tx)
	fmt.Println("#1")
	fmt.Printf("%# v\n", pretty.Formatter(invoiceItem))
	return invoiceItem
} 

func (i *Invoice) AddItemWithPrice(
	tx *gorm.DB, item *Item, 
	itemPrice *ItemPrice,
) *InvoiceItem {
	return i.AddInvoiceItem(tx, &InvoiceItem{
		InvoiceID: i.ID, ItemID: item.ID,
		ItemPriceID: itemPrice.ID,
	})
}	

func (i *Invoice) AddItem(
	tx *gorm.DB, item *Item,
) *InvoiceItem {
	return i.AddItemWithPrice(
		tx, item, item.LoadPrices(tx).CurrentPrice(),
	)
}

func (i *Invoice) AddItemByUUID(tx *gorm.DB, uuid uuid.UUID) *InvoiceItem {
	return i.AddItem(tx, FetchItem(tx, uuid))
}

func (i *Invoice) Calculate(tx *gorm.DB, loadItems ...bool) {
	var subTotal, taxable int

	// Auto load items
	if len(loadItems) < 1 || (len(loadItems) > 0 && loadItems[0] != false) {
		i.LoadItems(tx)
	}

	for _, item := range i.Items {
		subTotal += item.ItemPrice.Price
		
		if item.ItemPrice.Taxable == true {
			taxable += item.ItemPrice.Price
		}
	}

	// DemonDin Cut
	i.SubTotal = subTotal
	i.Taxes = int(float32(taxable) * 0.00)
	i.DemonDin = int(float32(subTotal) * 0.03)
	
	 // Calculate the Totals (Ensure ApplicationFee accuracy)
	// || \\ ^ [ i.Total = i.SubTotal + i.Taxes ]
	i.Total = i.SubTotal + i.Taxes + i.DemonDin
}

func (i *Invoice) Input(tx *gorm.DB, input *NewInvoice) {
	// Set Source Token
	if input.StripeTokenID != nil && *input.StripeTokenID != "" {
		i.SetStripeTokenID(tx, *input.StripeTokenID)
	}
	
	for _, _item := range input.Items {
		var invoiceItem *InvoiceItem

		// Create New InvoiceItem, or 'Fetch Existing InvoiceItem
		if _item.ID == nil || *_item.ID == uuid.Nil {
			invoiceItem = i.AddItemByUUID(tx, _item.ItemID)
		} else {
			invoiceItem = FetchInvoiceItem(tx, *_item.ID)
		}
		
		utils.PipeInput(_item, invoiceItem)
		i.AddInvoiceItem(tx, invoiceItem)
	}
}

func (i *Invoice) Save(tx *gorm.DB) {
	tx.Save(i)
}

func (i *Invoice) SetStripeTokenID(tx *gorm.DB, id string) {
	stripe.Key = os.Getenv("STRIPE_PRIVATE_KEY")
	tkn, err := token.Get(id, &stripe.TokenParams{})
	if err != nil {
		panic(err)
	}
	
	i.SetStripeToken(tx, tkn)
}

func (i *Invoice) SetStripeToken(tx *gorm.DB, tkn *stripe.Token) {
	i.StripeTokenID = stripe.String(tkn.ID)
	i.StripeToken = tkn
	i.Save(tx)
}

func (i *Invoice) Submit(tx *gorm.DB) error {
	if i.StripeTokenID == nil || *(i.StripeTokenID) == "" {
		return fmt.Errorf("Missing SourceToken at Invoice Submission")
	}

	 // List Items and Prices
	// || \\ ^ Update this for refunds!
	var items = map[string][]string{}
	for _, _item_ := range i.Items {
		var key = _item_.Sample()
		var opt = _item_.SampleOptions()
		items[key] = opt
	}

	var description = map[string]interface{}{} 
	description["subTotal"] = i.SubTotal
	description["demonDin"] = i.DemonDin
	description["taxes"] = i.Taxes
	description["total"] = i.Total
	description["items"] = items

	_description_, err := yaml.Marshal(description)
	if err != nil {
		return err
	}

	 // Stripe Transaction
	// || \\ ^
	stripe.Key = os.Getenv("STRIPE_PRIVATE_KEY")
	chargeParams := &stripe.ChargeParams{
		Amount: stripe.Int64(int64(i.SubTotal)),
		//ApplicationFeeAmount: stripe.Int64(int64(i.DemonDin)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String(string(_description_)),
	}
	
	//chargeParams.SetStripeAccount(os.Getenv("STRIPE_CONNECT_ACT"))
	chargeParams.SetSource(*(i.StripeTokenID)) // obtained with Stripe.js
	
	_charge_, err := charge.New(chargeParams)
	if err != nil {
		return err
	}

	i.StripeChargeID = stripe.String(_charge_.ID)
	i.StripeCharge = _charge_
	
	i.Save(tx)
	
	//fmt.Printf("%# v\n", pretty.Formatter(i))

	return nil
}

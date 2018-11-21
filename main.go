package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"gopkg.in/macaron.v1"

	"github.com/KellyLSB/demondin/models"
	stripeClient "github.com/stripe/stripe-go/client"
)

func main() {
	// Load .env file from repo
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load a macaron daemon
	m := macaron.Classic()
	// Injected automatically by macaron.Classic()
	//m.Use(macaron.Logger())
	//m.Use(macaron.Recovery())
	//m.Use(macaron.Static("public"))

	// Handle Database Connections
	m.Use(dbInit(
		os.Getenv("POSTGRES_HOSTPORT"), os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"),
	))

	// Handle StripeAPI
	m.Map(stripeClient.New(os.Getenv("STRIPE_PRIVATE_KEY"), nil))

	// Handle Templating
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Extensions: []string{".tmpl", ".html"},
		Directory:  "templates/default",
		IndentXML:  true,
		IndentJSON: true,

		Funcs: []template.FuncMap{map[string]interface{}{
			"AppName": func() string {
				return "DemonDin"
			},
			"AppVer": func() string {
				return "0.0.1"
			},
		}},
	}))

	// Routes
	m.Group("/shop", func() {
		m.Group("/keeper", func() {
			m.Get("/", func(
				ctx *macaron.Context,
				stripe *stripeClient.API,
			) {
				balance, _ := stripe.Balance.Get(nil)
				ctx.JSON(200, balance)
			})

			m.Get("/charges", func(
				ctx *macaron.Context,
				stripe *stripeClient.API,
			) {
				charges := stripe.Charges.List(nil)
				ctx.JSON(200, charges)
			})

			m.Get("/items?:ext(\\.[\\w]+$)", func(
			  ctx *macaron.Context,
			  db *gorm.DB,
			) {
				var badges []*models.Item

				db.Table("items").Select("*").Preload("Prices").Find(&badges)

				switch ctx.Params(":ext") {
				case ".json":
					ctx.JSON(200, badges)
				default:
					ctx.Data["Badges"] = badges
					ctx.Data["URL"] = ctx.URLFor("post_badge_price", ":item_id", "f31ac18e-0bb1-428e-952f-75cbc5604f3b", ":ext", ".json")
					ctx.HTML(200, "shop/keeper/badges")
				}
			})
			
			// I don't like the *.* (:path.:ext) routes option;
			// I'm feeling a little cornered here on that one
			m.Put("/items/:id([\\w-]+)?:ext(\\.[\\w]+$)", func(
				ctx *macaron.Context,
				db *gorm.DB, log *log.Logger,
			) {
				var badge models.Item
				db.Table("items").Select("*").Preload("Prices").
					Find(&badge, "items.id = ?", ctx.Params(":id"))

				decoder := json.NewDecoder(ctx.Req.Body().ReadCloser())
				if err := decoder.Decode(&badge); err == io.EOF {
					panic(err)
				} else if err != nil {
					panic(err)
				}

				db.Save(&badge)
				
				switch ctx.Params(":ext") {
				case ".json":
					ctx.JSON(200, badge)
				default:
					ctx.Data["Badge"] = badge
					ctx.HTML(200, "shop/keeper/badge")
				}
			})
			
			m.Post("/items/:item_id([\\w-]+)/prices?:ext(\\.[\\w]+$)", func(
				ctx *macaron.Context,
				db *gorm.DB, log *log.Logger,
			) {
				var price models.Price
				
				decoder := json.NewDecoder(ctx.Req.Body().ReadCloser())
				if err := decoder.Decode(&price); err == io.EOF {
					panic(err)
				} else if err != nil {
					panic(err)
				}
				
				// Set the ItemID
				price.ItemID = uuid.MustParse(ctx.Params(":item_id"))
				
				log.Println("Created new Record:\n%+v\n", price)
				db.Save(&price)
			}).Name("post_badge_price")
			
			m.Put("/items/:item([\\w-]+)/prices/:id([\\w-]+)?:ext(\\.[\\w]+$)", func(
				ctx *macaron.Context,
				db *gorm.DB, log *log.Logger,
			) {
				var price models.Price
				db.Table("prices").Select("*").
					Where("prices.item_id = ?", ctx.Params(":item")).
					Find(&price, "prices.id = ?", ctx.Params(":id"))
					
				decoder := json.NewDecoder(ctx.Req.Body().ReadCloser())
				if err := decoder.Decode(&price); err == io.EOF {
					panic(err)
				} else if err != nil {
					panic(err)
				}
				
				log.Println("Updated Record:\n%+v\n", price)
				db.Save(&price)
				
				switch ctx.Params(":ext") {
				case ".json":
					ctx.JSON(200, price)
				default:
					ctx.Data["Price"] = price
					ctx.HTML(200, "shop/keeper/price")
				}
			})
		})

		m.Get("/items?:ext(\\.[\\w]+$)", func(
		  ctx *macaron.Context,
		  db *gorm.DB, log *log.Logger,
		) {
			var badges []*models.Item

			db.Table("items").Select("*").
				//Joins("LEFT JOIN prices ON items.id = prices.item_id").
				Where("items.is_badge AND items.enabled").Preload("Prices",
				"? BETWEEN prices.valid_after AND prices.valid_before",
				time.Now(),
			).Find(&badges)

			switch ctx.Params(":ext") {
			case ".json":
				ctx.JSON(200, badges)
			default:
				// shop/badges.html template
				ctx.HTML(200, "shop/badges")
			}
		})
	})

	m.Get("/", myHandler)

	// Start our macaron daemon
	log.Println("Server is running...")
	log.Println(http.ListenAndServe(os.Getenv("HOSTPORT"), m))
}

func myHandler(ctx *macaron.Context) string {
	return "the request path is: " + ctx.Req.RequestURI
}

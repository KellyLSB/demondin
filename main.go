package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-macaron/pongo2"
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
	m.Use(pongo2.Pongoer())
	m.Map(stripeClient.New(os.Getenv("STRIPE_PRIVATE_KEY"), nil))
	m.Use(dbInit(
		os.Getenv("POSTGRES_HOSTPORT"), os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"),
	))

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
		})

		m.Get("/badges", func(ctx *macaron.Context, db *gorm.DB) {
			var badges []*models.Item

			db.Table("items").Select("*").
				//Joins("LEFT JOIN prices ON items.id = prices.item_id").
				Where("items.is_badge AND items.enabled").Preload("Prices",
				"? BETWEEN prices.valid_after AND prices.valid_before",
				time.Now(),
			).Find(&badges)

			ctx.JSON(200, badges)
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

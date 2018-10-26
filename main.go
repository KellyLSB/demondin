package main

import (
	"fmt"
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

	// Connect database
	db, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_NAME"),
		os.Getenv("POSTGRES_PASS"),
	))

	if err != nil {
		log.Fatalf("DB Connection Failed: %s", err)
	}

	defer db.Close()

	// Load a macaron daemon
	m := macaron.Classic()
	// Injected automatically by macaron.Classic()
	//m.Use(macaron.Logger())
	//m.Use(macaron.Recovery())
	//m.Use(macaron.Static("public"))
	m.Use(pongo2.Pongoer())
	m.Map(db)
	m.Map(stripeClient.New(os.Getenv("STRIPE_PRIVATE_KEY"), nil))

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
				Joins("LEFT JOIN prices ON items.id = prices.item_id").Where(
				"prices.valid_before <= ? && prices.valid_after > ? && "+
					"&& items.badge = ?", time.Now(), time.Now(), true,
			).Scan(&badges)

			ctx.JSON(200, badges)
		})
	})

	m.Get("/", myHandler)

	// Start our macaron daemon
	log.Println("Server is running...")
	log.Println(http.ListenAndServe(os.Getenv("HOST_PORT"), m))
}

func myHandler(ctx *macaron.Context) string {
	return "the request path is: " + ctx.Req.RequestURI
}

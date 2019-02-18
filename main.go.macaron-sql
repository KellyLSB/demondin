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
  "github.com/go-macaron/session"

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
	
	// Handle Session Cookies
	m.Use(session.Sessioner(session.Options{
    // Name of provider. Default is "memory".
    Provider:       "memory",
    // Provider configuration, it's corresponding to provider.
    ProviderConfig: "",
    // Cookie name to save session ID. Default is "MacaronSession".
    CookieName:     "demondin",
    // Cookie path to store. Default is "/".
    CookiePath:     "/",
    // GC interval time in seconds. Default is 3600.
    Gclifetime:     3600,
    // Max life time in seconds. Default is whatever GC interval time is.
    Maxlifetime:    3600,
    // Use HTTPS only. Default is false.
    Secure:         false,
    // Cookie life time. Default is 0.
    CookieLifeTime: 0,
    // Cookie domain name. Default is empty.
    Domain:         "",
    // Session ID length. Default is 16.
    IDLength:       16,
    // Configuration section name. Default is "session".
    Section:        "session",
}))

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
					ctx.HTML(200, "shop/keeper/items")
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
				
				//JSONDecoder(ctx, &price)
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
		  log *log.Logger, db *gorm.DB,
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
				ctx.HTML(200, "shop/items")
			}
		})
		
		m.Group("/invoicing", func() {
		  m.Get(".json", func(
		    ctx *macaron.Context, log *log.Logger,
		    sess session.Store, db *gorm.DB,
		  ) {
  		  invoicing := new(Invoicing)
  		  invoicing.LoadSession(sess, db, log)
  		  invoicing.SaveSession(sess, db, log)
  		  ctx.JSON(200, invoicing.Invoice)
		  })
		
  		m.Post("/addToCart.json", func(
  		  ctx *macaron.Context, log *log.Logger,
  		  sess session.Store, db *gorm.DB,
  		) {
  		  invoicing := new(Invoicing)
  		  invoicing.LoadSession(sess, db, log)
  		  invoicing.PostAddToCart(ctx, db, log)
  		  invoicing.SaveSession(sess, db, log)
  		  ctx.JSON(200, invoicing.Invoice)
  		})
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

func JSONDecoder(ctx *macaron.Context, into interface{}) {
  decoder := json.NewDecoder(ctx.Req.Body().ReadCloser())
  if err := decoder.Decode(into); err == io.EOF {
	  panic(err)
  } else if err != nil {
	  panic(err)
  }
}


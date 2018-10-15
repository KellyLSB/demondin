package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"

	"github.com/KellyLSB/youmacon/app/controllers"
	"github.com/KellyLSB/youmacon/app/controllers/shop"
	"github.com/KellyLSB/youmacon/app/controllers/shop/keeper"

	stripeClient "github.com/stripe/stripe-go/client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := iris.New()

	// Define templates using the std html/template engine.
	// Parse and load all files inside "./views" folder with ".html" file extension.
	// Reload the templates on each request (development mode).
	app.RegisterView(iris.HTML("./app/views", ".html").Reload(true))

	// Log request
	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		ctx.Next()
	})

	// Register custom handler for specific http errors.
	app.OnAnyErrorCode(func(ctx iris.Context) {
		// .Values are used to communicate between handlers, middleware.
		errMessage := ctx.Values().GetString("error")

		if errMessage != "" {
			ctx.Writef("Internal server error: %s", errMessage)
			ctx.Application().Logger().Errorf("Internal server error: %s", errMessage)
			return
		}

		ctx.Writef("(Unexpected) internal server error")
		ctx.Application().Logger().Errorf("(Unexpected) internal server error")
	})

	mvc.New(app.Party("/")).Register(func(ctx iris.Context) (DB *gorm.DB) {
		db, err := gorm.Open("postgres", fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s",
			os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_NAME"),
			os.Getenv("POSTGRES_PASS"),
		))

		db.LogMode(true)
		db.SetLogger(ctx.Application().Logger())

		if err != nil {
			ctx.Application().Logger().Errorf("DB Connection Failed: %s", err)
		}

		dbMigrate(db)

		iris.RegisterOnInterrupt(func() {
			db.Close()
		})

		return db
	}).Configure(func(mvcApp *mvc.Application) {
		mvcApp.Handle(new(controllers.Index))
		mvcApp.Party("/shop").Register(stripeService).Handle(new(shop.Index))
		mvcApp.Party("/shop/keeper").Register(stripeService).Handle(new(shopkeeper.Index))
	})

	app.Run(iris.Addr(os.Getenv("HOST_PORT")), iris.WithCharset("UTF-8"))
}

func stripeService(ctx iris.Context) (StripeAPI *stripeClient.API) {
	return stripeClient.New(os.Getenv("STRIPE_KEY"), nil)
}

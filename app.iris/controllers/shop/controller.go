package shop

import (
	"github.com/KellyLSB/youmacon/app/controllers"
	"github.com/stripe/stripe-go/client"
)

type Controller struct {
	controllers.Controller
	// StripeKey pull from an env variable later.
	StripeAPI *client.API
}

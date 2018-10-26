package shopkeeper

import (
	"github.com/kataras/iris"
)

type Index struct {
	Controller
}

func (i Index) Get(ctx iris.Context) {
	balance, err := i.StripeAPI.Balance.Get(nil)
	i.Panic(err)
	ctx.JSON(balance)
}

func (i Index) GetCharges(ctx iris.Context) {
	charges := i.StripeAPI.Charges.List(nil)
	ctx.JSON(charges)
}

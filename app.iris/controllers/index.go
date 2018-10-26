package controllers

import "github.com/kataras/iris"

// Controller is the inherited http controller structure.
type Index struct {
	Controller
}

func (i Index) Get(ctx iris.Context) {
	ctx.JSON(map[string]interface{}{
		"pizza": "pie",
	})
}

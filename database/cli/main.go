package main

import (
	"github.com/KellyLSB/demondin/database"
	"github.com/KellyLSB/demondin/graphql/model"
	"github.com/kr/pretty"
	"fmt"
)

func main() {
	fmt.Printf("%# v\n", pretty.Formatter(database.LoadModelStruct(
		&model.Invoice{},
	).GetFields()))
}

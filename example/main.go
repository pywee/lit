package main

import (
	"fmt"
	"os"

	goExpr "github.com/pywee/lit"
)

func main() {
	exprs, err := os.ReadFile("./testing.lit")
	if err != nil {
		panic(err)
	}
	if _, err = goExpr.NewExpr(exprs); err != nil {
		fmt.Println(err)
		return
	}
}

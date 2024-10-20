package main

import (
	"fmt"

	"github.com/david22573/gnotes/internal/router"
)

func main() {
	fmt.Println("Hi")
	r := router.NewRouter()
	r.Logger.Fatal(r.Start(":1323"))
}

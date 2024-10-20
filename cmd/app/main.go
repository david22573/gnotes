package main

import (
	"github.com/david22573/gnotes/internal/server"
)

func main() {
	s := server.Init()
	s.Logger.Fatal(s.Start(":1323"))
}

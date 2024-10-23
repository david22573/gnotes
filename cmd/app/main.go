package main

import (
	"log"

	"github.com/david22573/gnotes/internal/server"
)

func main() {
	srv := server.New("data.db")
	if err := srv.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

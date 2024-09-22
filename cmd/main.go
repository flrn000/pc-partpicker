package main

import (
	"log"

	"github.com/flrn000/pc-partpicker/cmd/api"
)

func main() {
	server := api.NewAPIServer(":8080", nil)
	if err := server.Start(); err != nil {
		log.Fatalf("starting server: %v", err)
	}
}

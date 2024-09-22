package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/flrn000/pc-partpicker/cmd/api"
	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbpool, err := db.NewPSQLStorage(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	defer dbpool.Close()

	// Ping the database to verify the connection
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	fmt.Println("Connected to the database successfully!")

	userStore := data.NewUserStore(dbpool)

	server := api.NewAPIServer(":8080", userStore)
	if err := server.Start(); err != nil {
		log.Fatalf("starting server: %v", err)
	}
}

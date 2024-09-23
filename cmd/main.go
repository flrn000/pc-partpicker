package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/flrn000/pc-partpicker/cmd/api"
	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/db"
	"github.com/flrn000/pc-partpicker/logging"
	"github.com/joho/godotenv"
)

func main() {
	logger := logging.NewLogger(slog.NewTextHandler(os.Stderr, nil))

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbpool, err := db.NewPSQLStorage(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to create connection pool: %v\n", err))
		os.Exit(1)
	}

	addr := flag.String("address", ":8080", "HTTP network address")

	flag.Parse()

	defer dbpool.Close()

	// Ping the database to verify the connection
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	logger.Info("Connected to the database successfully!")

	userStore := data.NewUserStore(dbpool)

	server := api.NewAPIServer(*addr, logger, userStore)
	if err := server.Start(); err != nil {
		log.Fatalf("starting server: %v", err)
	}
}

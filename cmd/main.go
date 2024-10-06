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
	"github.com/flrn000/pc-partpicker/types"
	"github.com/joho/godotenv"
)

func main() {
	appConfig := &types.AppConfig{
		Logger: logging.NewLogger(slog.NewTextHandler(os.Stderr, nil)),
	}

	err := godotenv.Load()
	if err != nil {
		appConfig.Logger.Error("Error loading .env file")
		os.Exit(1)
	}

	flag.StringVar(&appConfig.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&appConfig.Address, "address", ":8080", "HTTP network address")

	flag.StringVar(&appConfig.DB_URL, "db-dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")

	flag.StringVar(&appConfig.JWTSecret, "jwtSecret", os.Getenv("JWT_SECRET"), "JWT Secret")

	dbpool, err := db.NewPSQLStorage(appConfig.DB_URL)
	if err != nil {
		appConfig.Logger.Error(fmt.Sprintf("Unable to create connection pool: %v\n", err))
		os.Exit(1)
	}

	flag.Parse()

	defer dbpool.Close()

	// Ping the database to verify the connection
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	appConfig.Logger.Info("Connected to the database successfully!")

	userStore := data.NewUserStore(dbpool)
	componentStore := data.NewComponentStore(dbpool)
	refreshTokenStore := data.NewRefreshTokenStore(dbpool)

	server := api.NewAPIServer(
		appConfig,
		userStore,
		refreshTokenStore,
		componentStore,
	)
	if err := server.Start(); err != nil {
		appConfig.Logger.Error(fmt.Sprintf("starting server: %v", err))
		os.Exit(1)
	}
}

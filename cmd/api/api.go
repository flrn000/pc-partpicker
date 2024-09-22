package api

import (
	"log"
	"net/http"
	"time"

	"github.com/flrn000/pc-partpicker/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServer struct {
	address string
	db      *pgxpool.Pool
}

func NewAPIServer(addr string, db *pgxpool.Pool) *APIServer {
	return &APIServer{
		address: addr,
		db:      db,
	}
}

func (s *APIServer) Start() error {
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler:      mux,
		Addr:         s.address,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	service.AddRoutes(mux)

	log.Println("Listening on", s.address)

	return srv.ListenAndServe()
}

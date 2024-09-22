package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/flrn000/pc-partpicker/service"
)

type APIServer struct {
	address string
	db      *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
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

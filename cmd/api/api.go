package api

import (
	"log"
	"net/http"
	"time"

	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/service"
)

type APIServer struct {
	address   string
	userStore *data.UserStore
}

func NewAPIServer(addr string, userStore *data.UserStore) *APIServer {
	return &APIServer{
		address:   addr,
		userStore: userStore,
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

	service.AddRoutes(mux, s.userStore)

	log.Println("Listening on", s.address)

	return srv.ListenAndServe()
}

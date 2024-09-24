package api

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/middleware"
	"github.com/flrn000/pc-partpicker/service"
)

type APIServer struct {
	address        string
	logger         *slog.Logger
	userStore      *data.UserStore
	componentStore *data.ComponentStore
}

func NewAPIServer(
	addr string,
	logger *slog.Logger,
	userStore *data.UserStore,
	componentStore *data.ComponentStore,
) *APIServer {
	return &APIServer{
		address:        addr,
		logger:         logger,
		userStore:      userStore,
		componentStore: componentStore,
	}
}

func (s *APIServer) Start() error {
	mux := http.NewServeMux()
	var handler http.Handler = mux
	addLogging := middleware.NewLogging(s.logger)
	handler = middleware.AddSecureHeaders(addLogging(handler))

	srv := &http.Server{
		Handler:      handler,
		Addr:         s.address,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	service.AddRoutes(mux, s.userStore, s.componentStore)

	log.Println("Listening on", s.address)

	return srv.ListenAndServe()
}

package api

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/middleware"
	"github.com/flrn000/pc-partpicker/service"
	"github.com/flrn000/pc-partpicker/utils"
)

type APIServer struct {
	address        string
	validator      *utils.Validator
	logger         *slog.Logger
	userStore      *data.UserStore
	componentStore *data.ComponentStore
}

func NewAPIServer(
	addr string,
	validator *utils.Validator,
	logger *slog.Logger,
	userStore *data.UserStore,
	componentStore *data.ComponentStore,
) *APIServer {
	return &APIServer{
		address:        addr,
		validator:      validator,
		logger:         logger,
		userStore:      userStore,
		componentStore: componentStore,
	}
}

func (s *APIServer) Start() error {
	mux := http.NewServeMux()
	var handler http.Handler = mux
	addLogging := middleware.NewLogging(s.logger)
	handler = addLogging(handler)
	handler = middleware.AddSecureHeaders(handler)

	srv := &http.Server{
		Handler:      handler,
		Addr:         s.address,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	service.AddRoutes(mux, s.validator, s.userStore, s.componentStore)

	log.Println("Listening on", s.address)

	return srv.ListenAndServe()
}

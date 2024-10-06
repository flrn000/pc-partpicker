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
	address           string
	jwtSecret         string
	logger            *slog.Logger
	userStore         *data.UserStore
	refreshTokenStore *data.RefreshTokenStore
	componentStore    *data.ComponentStore
}

func NewAPIServer(
	addr string,
	jwtSecret string,
	logger *slog.Logger,
	userStore *data.UserStore,
	refreshTokenStore *data.RefreshTokenStore,
	componentStore *data.ComponentStore,
) *APIServer {
	return &APIServer{
		address:           addr,
		jwtSecret:         jwtSecret,
		logger:            logger,
		userStore:         userStore,
		refreshTokenStore: refreshTokenStore,
		componentStore:    componentStore,
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
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	service.AddRoutes(
		mux,
		s.jwtSecret,
		s.userStore,
		s.refreshTokenStore,
		s.componentStore,
	)

	log.Println("Listening on", s.address)

	return srv.ListenAndServe()
}

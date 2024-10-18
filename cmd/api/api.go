package api

import (
	"log"
	"net/http"
	"time"

	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/middleware"
	"github.com/flrn000/pc-partpicker/service"
	"github.com/flrn000/pc-partpicker/types"
)

type APIServer struct {
	appConfig         *types.AppConfig
	userStore         *data.UserStore
	refreshTokenStore *data.RefreshTokenStore
	componentStore    *data.ComponentStore
}

func NewAPIServer(
	appConfig *types.AppConfig,
	userStore *data.UserStore,
	refreshTokenStore *data.RefreshTokenStore,
	componentStore *data.ComponentStore,
) *APIServer {
	return &APIServer{
		appConfig:         appConfig,
		userStore:         userStore,
		refreshTokenStore: refreshTokenStore,
		componentStore:    componentStore,
	}
}

func (s *APIServer) Start() error {
	mux := http.NewServeMux()
	var handler http.Handler = mux

	handler = middleware.RecoverPanic(handler)
	addLogging := middleware.NewLogging(s.appConfig.Logger)
	handler = addLogging(handler)
	handler = middleware.AddSecureHeaders(handler)

	srv := &http.Server{
		Handler:      handler,
		Addr:         s.appConfig.Address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	service.AddRoutes(
		mux,
		s.appConfig,
		s.userStore,
		s.refreshTokenStore,
		s.componentStore,
	)

	log.Println("Listening on", s.appConfig.Address)

	return srv.ListenAndServe()
}

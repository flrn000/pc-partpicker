package service

import (
	"net/http"

	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/middleware"
	"github.com/flrn000/pc-partpicker/types"
)

func AddRoutes(
	mux *http.ServeMux,
	appConfig *types.AppConfig,
	userStore *data.UserStore,
	refreshTokenStore *data.RefreshTokenStore,
	componentStore *data.ComponentStore,
) {
	authenticate := middleware.WithAuthenticate(appConfig.JWTSecret, userStore)

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./client/static"))))

	mux.Handle("GET /{$}", middleware.RateLimit(handleIndex()))
	mux.Handle("GET /accounts/register", middleware.RateLimit(handleRegisterPage()))
	mux.Handle("GET /accounts/login", middleware.RateLimit(handleLoginPage()))
	mux.Handle("GET /products/{componentType}", middleware.RateLimit(handleViewProducts(componentStore)))

	mux.Handle("POST /api/v1/login", middleware.RateLimit(handleLogin(userStore, refreshTokenStore, appConfig.JWTSecret)))
	mux.Handle("POST /api/v1/register", middleware.RateLimit(handleRegister(userStore)))

	mux.Handle("POST /api/v1/products", middleware.RateLimit(authenticate(handleCreateProducts(componentStore))))
	mux.Handle("POST /api/v1/refresh", middleware.RateLimit(handleRefresh(refreshTokenStore, appConfig.JWTSecret)))
}

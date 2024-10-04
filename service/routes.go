package service

import (
	"net/http"

	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/middleware"
)

func AddRoutes(
	mux *http.ServeMux,
	jwtSecret string,
	userStore *data.UserStore,
	componentStore *data.ComponentStore,
) {
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./client/static"))))

	mux.Handle("GET /{$}", middleware.RateLimit(handleIndex()))
	mux.Handle("GET /accounts/register", middleware.RateLimit(handleRegisterPage()))
	mux.Handle("GET /accounts/login", middleware.RateLimit(handleLoginPage()))
	mux.Handle("GET /products/{componentType}", middleware.RateLimit(handleViewProducts(componentStore)))

	mux.Handle("POST /api/v1/login", middleware.RateLimit(handleLogin(userStore, jwtSecret)))
	mux.Handle("POST /api/v1/register", middleware.RateLimit(handleRegister(userStore)))

	mux.Handle("POST /api/v1/products", middleware.RateLimit(handleCreateProducts(componentStore)))
}

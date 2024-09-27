package service

import (
	"net/http"

	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/utils"
)

func AddRoutes(
	mux *http.ServeMux,
	validator *utils.Validator,
	userStore *data.UserStore,
	componentStore *data.ComponentStore,
) {
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./client/static"))))

	mux.Handle("GET /{$}", handleIndex())
	mux.Handle("GET /accounts/register", handleRegisterPage())
	mux.Handle("GET /products/{componentType}", handleViewProducts(componentStore))

	mux.Handle("POST /api/v1/login", handleLogin(userStore))
	mux.Handle("POST /api/v1/register", handleRegister(userStore, validator))

	mux.Handle("POST /api/v1/products", handleCreateProducts(componentStore))
}

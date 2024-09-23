package service

import (
	"net/http"

	"github.com/flrn000/pc-partpicker/data"
)

func AddRoutes(mux *http.ServeMux, userStore *data.UserStore) {
	mux.Handle("GET /{$}", handleIndex())
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./client/static"))))

	mux.Handle("POST /api/v1/login", handleLogin(userStore))
	mux.Handle("POST /api/v1/register", handleRegister(userStore))
}

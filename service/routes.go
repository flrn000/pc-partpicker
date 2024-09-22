package service

import (
	"net/http"

	"github.com/flrn000/pc-partpicker/data"
)

func AddRoutes(mux *http.ServeMux, userStore *data.UserStore) {
	mux.Handle("GET /login", handleLogin(userStore))
	mux.Handle("POST /api/v1/register", handleRegister(userStore))
}

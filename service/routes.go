package service

import "net/http"

func AddRoutes(mux *http.ServeMux) {
	mux.Handle("GET /login", handleLogin())
	mux.Handle("POST /register", handleRegister())
}

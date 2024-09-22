package service

import "net/http"

func handleLogin() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

		},
	)
}

func handleRegister() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

		},
	)
}

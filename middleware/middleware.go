package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/types"
	"github.com/flrn000/pc-partpicker/utils"
	"golang.org/x/time/rate"
)

func NewLogging(logger *slog.Logger) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				logger.Info("Incoming request", "method", r.Method, "url", r.URL.Path)

				next.ServeHTTP(w, r)
			},
		)
	}
}

func AddSecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Security-Policy",
				"default-src 'self'; script-src 'self' cdn.jsdelivr.net; style-src 'self' 'unsafe-inline' fonts.googleapis.com cdn.jsdelivr.net; font-src fonts.gstatic.com")

			w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "deny")
			w.Header().Set("X-XSS-Protection", "0")

			next.ServeHTTP(w, r)
		},
	)
}

func RateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			utils.WriteError(w, r, http.StatusTooManyRequests, errors.New("too many requests"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func WithAuthenticate(jwtSecret string, userStore *data.UserStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				token, err := utils.GetAuthToken(r.Header)
				if err != nil {
					w.Header().Set("WWW-Authenticate", "Bearer")
					utils.WriteError(w, r, http.StatusUnauthorized, err)
					return
				}

				userID, err := utils.ValidateJWT(token, jwtSecret)
				if err != nil {
					utils.WriteError(w, r, http.StatusUnauthorized, err)
					return
				}

				r = utils.ContextSetUser(r, &types.User{ID: userID})

				next.ServeHTTP(w, r)
			},
		)
	}
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// This will always be run in case a panic occurs as go unwinds the stack
			defer func() {
				if err := recover(); err != nil {
					// w.Header().Set("Connection", "close")
					utils.WriteError(w, r, http.StatusInternalServerError, fmt.Errorf("%v", err))
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}

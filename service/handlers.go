package service

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/types"
	"github.com/flrn000/pc-partpicker/utils"
	"golang.org/x/crypto/bcrypt"
)

func handleIndex() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			files := []string{
				"./client/html/base.tmpl",
				"./client/html/partials/nav.tmpl",
				"./client/html/pages/home.tmpl",
			}
			ts, err := template.ParseFiles(files...)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}

			err = ts.ExecuteTemplate(w, "base", nil)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}
		},
	)
}

func handleLogin(userStore *data.UserStore) http.Handler {
	type LoginUserPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

		},
	)
}

func handleRegister(userStore *data.UserStore) http.Handler {
	type RegisterUserPayload struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			val, err := utils.Decode[RegisterUserPayload](r)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}

			if _, err := userStore.GetByEmail(val.Email); err == nil {
				utils.WriteError(w, r, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", val.Email))
				return
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(val.Password), bcrypt.DefaultCost)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}

			user := &types.User{
				UserName: val.UserName,
				Email:    val.Email,
				Password: string(hashedPassword),
			}

			err = userStore.Create(user)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}

			w.Header().Set("Location", fmt.Sprintf("api/v1/users/%d", user.ID))

			err = utils.Encode(w, r, http.StatusCreated, user)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}
		},
	)
}

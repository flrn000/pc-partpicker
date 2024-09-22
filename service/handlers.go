package service

import (
	"fmt"
	"net/http"

	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/types"
	"github.com/flrn000/pc-partpicker/utils"
	"golang.org/x/crypto/bcrypt"
)

func handleLogin(userStore *data.UserStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

		},
	)
}

func handleRegister(userStore *data.UserStore) http.Handler {
	type RegisterUserPayload struct {
		UserName string `json:"user_name"`
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

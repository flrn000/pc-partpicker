package service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/flrn000/pc-partpicker/data"
	"github.com/flrn000/pc-partpicker/types"
	"github.com/flrn000/pc-partpicker/utils"
	"golang.org/x/crypto/bcrypt"
)

type templateData struct {
	User          *types.User
	Component     *types.Component
	ComponentType string
	Components    []*types.Component
	Form          any
}

func handleIndex() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			utils.RenderTemplate(w, r, http.StatusOK, "home.tmpl", nil)
		},
	)
}

func handleRegisterPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			utils.RenderTemplate(w, r, http.StatusOK, "register.tmpl", nil)
		},
	)
}

func handleLogin(userStore *data.UserStore, jwtSecret string) http.Handler {
	type LoginUserPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type ResponsePayload struct {
		ID    int    `json:"id"`
		Token string `json:"token"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			payload, err := utils.Decode[LoginUserPayload](r)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}
			u, err := userStore.GetByEmail(payload.Email)
			if err != nil {
				utils.WriteError(w, r, http.StatusUnauthorized, errors.New("email or password are incorrect"))
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(payload.Password))
			if err != nil {
				utils.WriteError(w, r, http.StatusUnauthorized, errors.New("email or password are incorrect"))
				return
			}

			token, err := utils.GenerateJWT(strconv.Itoa(u.ID), jwtSecret, time.Hour)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}

			result := ResponsePayload{
				ID:    u.ID,
				Token: token,
			}

			w.Header().Set("Location", "/")
			utils.Encode(w, r, http.StatusSeeOther, result)
		},
	)
}

func handleRegister(userStore *data.UserStore, validator *utils.Validator) http.Handler {
	type registerForm struct {
		Username    string
		Email       string
		FieldErrors map[string]string
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				utils.WriteError(w, r, http.StatusBadRequest, err)
				return
			}

			username := r.PostForm.Get("username")
			email := r.PostForm.Get("email")
			password := r.PostForm.Get("password")

			validator.CheckField(utils.IsNotBlank(username), "username", "The username cannot be blank")
			validator.CheckField(utils.IsNotBlank(email), "email", "The email cannot be blank")
			validator.CheckField(utils.IsNotBlank(password), "password", "The password cannot be blank")

			if !validator.IsValid() {
				data := templateData{
					Form: registerForm{
						Username:    username,
						Email:       email,
						FieldErrors: validator.FieldErrors,
					},
				}

				utils.RenderTemplate(w, r, http.StatusUnprocessableEntity, "register.tmpl", data)

				return
			}

			if _, err := userStore.GetByEmail(email); err == nil {
				utils.WriteError(w, r, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", email))
				return
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}

			user := &types.User{
				UserName: username,
				Email:    email,
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

func handleCreateProducts(componentStore *data.ComponentStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			payload, err := utils.Decode[types.CreateProductPayload](r)
			if err != nil {
				utils.WriteError(w, r, http.StatusBadRequest, err)
				return
			}

			component := &types.Component{
				Name:         payload.Name,
				Type:         payload.Type,
				Manufacturer: payload.Manufacturer,
				Model:        payload.Model,
				Price:        payload.Price,
				Rating:       payload.Rating,
				ImageURL:     payload.ImageURL,
			}

			err = componentStore.Create(component)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
			}

			w.Header().Set("Location", fmt.Sprintf("api/v1/products/%d", component.ID))

			err = utils.Encode(w, r, http.StatusCreated, component)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}
		})
}

func handleViewProducts(componentStore *data.ComponentStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			componentType := r.PathValue("componentType")

			components, err := componentStore.GetMany(60, types.ComponentType(componentType))
			if err != nil {
				utils.WriteError(w, r, http.StatusBadRequest, err)
				return
			}

			data := templateData{
				ComponentType: componentType,
				Components:    components,
			}

			utils.RenderTemplate(w, r, http.StatusOK, "products.tmpl", data)
		},
	)
}

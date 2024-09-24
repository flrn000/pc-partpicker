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
			payload, err := utils.Decode[LoginUserPayload](r)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}
			u, err := userStore.GetByEmail(payload.Email)
			if err != nil {
				utils.WriteError(w, r, http.StatusBadRequest, fmt.Errorf("user with email %s doesn't exist", payload.Email))
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(payload.Password))
			if err != nil {
				utils.WriteError(w, r, http.StatusUnauthorized, err)
				return
			}
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

			components, err := componentStore.GetMany(10, types.ComponentType(componentType))
			if err != nil {
				utils.WriteError(w, r, http.StatusBadRequest, err)
				return
			}

			templateData := struct {
				ComponentType string
				Components    []*types.Component
			}{
				componentType,
				components,
			}

			files := []string{
				"./client/html/base.tmpl",
				"./client/html/partials/nav.tmpl",
				"./client/html/pages/products.tmpl",
			}
			ts, err := template.ParseFiles(files...)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}

			err = ts.ExecuteTemplate(w, "base", templateData)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}
		},
	)
}

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

func handleLoginPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			utils.RenderTemplate(w, r, http.StatusOK, "login.tmpl", nil)
		},
	)
}

func handleLogin(userStore *data.UserStore, jwtSecret string) http.Handler {
	type loginForm struct {
		Email    string
		Password string
		utils.Validator
	}
	type ResponsePayload struct {
		ID    int    `json:"id"`
		Token string `json:"token"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}

			form := loginForm{
				Email:    r.PostForm.Get("email"),
				Password: r.PostForm.Get("password"),
			}

			form.CheckField(utils.IsNotBlank(form.Email), "email", "Email cannot be blank")
			form.CheckField(utils.IsValidEmail(form.Email), "email", "Email must be a valid address")
			form.CheckField(utils.IsNotBlank(form.Password), "password", "Password cannot be blank")

			if !form.IsValid() {
				data := templateData{
					Form: form,
				}

				utils.RenderTemplate(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)

				return
			}

			u, err := userStore.GetByEmail(form.Email)
			if err != nil {
				form.AddFieldError("email", "Email address is incorrect")
				data := templateData{
					Form: form,
				}

				utils.RenderTemplate(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(form.Password))
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

func handleRegister(userStore *data.UserStore) http.Handler {
	type registerForm struct {
		Username string
		Email    string
		Password string
		utils.Validator
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				utils.WriteError(w, r, http.StatusBadRequest, err)
				return
			}

			form := registerForm{
				Username: r.PostForm.Get("username"),
				Email:    r.PostForm.Get("email"),
				Password: r.PostForm.Get("password"),
			}

			form.CheckField(utils.IsNotBlank(form.Username), "username", "The username cannot be blank")
			form.CheckField(utils.IsNotBlank(form.Email), "email", "The email cannot be blank")
			form.CheckField(utils.IsValidEmail(form.Email), "email", "Email must be a valid address")
			form.CheckField(utils.IsNotBlank(form.Password), "password", "The password cannot be blank")
			form.CheckField(utils.MinChars(form.Password, 8), "password", "The password must be at least 8 characters long")

			if !form.IsValid() {
				data := templateData{
					Form: form,
				}

				utils.RenderTemplate(w, r, http.StatusUnprocessableEntity, "register.tmpl", data)

				return
			}

			if _, err := userStore.GetByEmail(form.Email); err == nil {
				form.AddFieldError("email", "Email address is already in use")
				data := templateData{
					Form: form,
				}

				utils.RenderTemplate(w, r, http.StatusUnprocessableEntity, "register.tmpl", data)
				return
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}

			user := &types.User{
				UserName: form.Username,
				Email:    form.Email,
				Password: string(hashedPassword),
			}

			err = userStore.Create(user)
			if err != nil {
				utils.WriteError(w, r, http.StatusInternalServerError, err)
				return
			}

			http.Redirect(w, r, "/accounts/login", http.StatusSeeOther)
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

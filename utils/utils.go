package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) IsValid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func IsNotBlank(val string) bool {
	return strings.TrimSpace(val) != ""
}

func MaxChars(val string, n int) bool {
	return utf8.RuneCountInString(val) <= n
}

func IsPermittedInt(val int, permittedValues ...int) bool {
	for _, v := range permittedValues {
		if val == v {
			return true
		}
	}
	return false
}

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var val T
	if err := json.NewDecoder(r.Body).Decode(&val); err != nil {
		return val, fmt.Errorf("decode json: %w", err)
	}
	return val, nil
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, err error) {
	slog.Error(err.Error())
	Encode(w, r, status, map[string]string{"error": err.Error()})
}

func RenderTemplate(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	page string,
	data any) {
	files := []string{
		"./client/html/base.tmpl",
		"./client/html/partials/nav.tmpl",
		fmt.Sprintf("./client/html/pages/%s", page),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	buf := new(bytes.Buffer)

	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

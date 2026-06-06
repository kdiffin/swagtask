package router

import (
	"net/http"
	"swagtask/internal/auth"
)

func SetupAuthRoutes(mux *http.ServeMux, handler *auth.AuthHandler) {
	mux.HandleFunc("POST /sign-up/{$}", handler.Signup)
	mux.HandleFunc("GET /sign-up/{$}", handler.SignupPage)
	mux.HandleFunc("POST /login/{$}", handler.Login)
	mux.HandleFunc("GET /login/{$}", handler.LoginPage)
	mux.HandleFunc("POST /logout/{$}", handler.Logout)
}

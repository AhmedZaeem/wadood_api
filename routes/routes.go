package routes

import (
	"net/http"
	"wadood/handlers"
)

func Setup() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("POST /register", handlers.RegisterHandler)
	r.HandleFunc("POST /login", handlers.LoginHandler)
	r.HandleFunc("POST /logout", handlers.LogoutHandler)
	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Welcome to Wadood API"))
	})
	r.HandleFunc("POST /profile/edit", handlers.EditProfileHandler)

	return r
}
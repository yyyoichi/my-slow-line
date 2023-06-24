package main

import (
	"fmt"
	"himakiwa/handlers"
	"himakiwa/handlers/middleware"
	"himakiwa/services/database"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func init() {
	database.Connect()
}

func main() {
	fmt.Println("hello")
	handler()
}

func handler() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/safe", Safe).Methods(http.MethodGet)

	ah := handlers.NewAutenticateHandlers()
	api.HandleFunc("/signin", ah.SigninHandler).Methods(http.MethodPost)
	api.HandleFunc("/login", ah.LoginHandler).Methods(http.MethodPost)
	api.HandleFunc("/codein", ah.VerificateHandler).Methods(http.MethodPost)

	me := api.PathPrefix("/me").Subrouter()
	me.HandleFunc("/", handlers.MeHandler).Methods(http.MethodGet)
	me.HandleFunc("/logout", handlers.LogoutHandler).Methods(http.MethodPost)
	me.Use(middleware.AuthMiddleware)

	wp := api.PathPrefix("/webpush").Subrouter()
	wp.HandleFunc("/vapid_public_key", handlers.VapidHandler).Methods(http.MethodGet)
	wp.HandleFunc("/subscription", handlers.PushSubscriptionHandler).Methods(http.MethodGet, http.MethodPost, http.MethodDelete)

	api.Use(middleware.CSRFMiddleware)
	api.Use(middleware.CROSMiddleware)
	http.ListenAndServe(":8080", r)
}

func Safe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	w.WriteHeader(http.StatusOK)
}

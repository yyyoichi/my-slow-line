package main

import (
	"encoding/json"
	"fmt"
	"himakiwa/database"
	"himakiwa/handlers"
	"himakiwa/middleware"
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
	api.HandleFunc("/", Index).Methods(http.MethodGet)
	api.HandleFunc("/safe", Safe).Methods(http.MethodGet)
	api.HandleFunc("/post", Post).Methods(http.MethodPost)
	api.HandleFunc("/signin", handlers.SigninHandler).Methods(http.MethodPost)
	api.HandleFunc("/me", handlers.MeHandler).Methods(http.MethodGet)
	api.Use(middleware.CROSMiddleware)
	api.Use(middleware.CSRFMiddleware)
	http.ListenAndServe(":8080", r)
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	p := Person{"yama", 23}
	json.NewEncoder(w).Encode(p)
}

func Safe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	w.WriteHeader(http.StatusOK)
}

func Post(w http.ResponseWriter, r *http.Request) {
	var p Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
	fmt.Printf("%s: %d \n", p.Name, p.Age)

	json.NewEncoder(w).Encode(p)
}

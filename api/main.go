package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("hello")
	handler()
}
func handler() {
	r := mux.NewRouter()
	r.HandleFunc("/", Index).Methods(http.MethodGet)
	r.HandleFunc("/safe", Safe).Methods(http.MethodGet)
	r.HandleFunc("/post", Post).Methods(http.MethodPost)
	r.HandleFunc("/post", Preflight).Methods(http.MethodOptions)
	r.Use(CROSMiddleware)
	csrfMiddleware := getCsrfMiddleware()
	r.Use(csrfMiddleware)
	http.ListenAndServe(":8080", r)
}
func CROSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Csrf-Token")
		w.Header().Set("Access-Control-Expose-Headers", "X-Csrf-Token")
		fmt.Printf("got from '%s' method '%s' to '%s'\n", r.Host, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func getCsrfMiddleware() func(http.Handler) http.Handler {
	key := []byte("xox-oxo-xox-oxo-xox-oxo-xox-oxo-")
	return csrf.Protect(key,
		csrf.Secure(false),
		csrf.HttpOnly(false),
		csrf.TrustedOrigins([]string{"http://127.0.0.1:3000", "http://localhost:3000"}),
		csrf.ErrorHandler(http.HandlerFunc(serverError)),
	)
}

func serverError(w http.ResponseWriter, r *http.Request) {
	fmt.Println("error csrf token")
	http.Error(w, "error", http.StatusBadRequest)
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
	return
}

func Preflight(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
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

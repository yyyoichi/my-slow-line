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
	r.HandleFunc("/", Index).Methods("GET")
	key := []byte("xox-oxo-xox-oxo-xox-oxo-xox-oxo-")
	http.ListenAndServe(":8000",
		csrf.Protect(key,
			csrf.Secure(false),
			csrf.HttpOnly(false),
			csrf.TrustedOrigins([]string{"http://localhost:3000"}),
		)(r))
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	p := Person{"yama", 23}
	json.NewEncoder(w).Encode(p)
}

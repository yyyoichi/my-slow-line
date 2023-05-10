package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var (
	InvalidEndPointQuery = "invalid query"
)

func VapidHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(os.Getenv("VAPID_PUBLIC_KEY"))
}

func PushSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getPushSubscriptionHandler(w, r)
	case "POST":
		postPushSubscriptionHandler(w, r)
	case "DELETE":
		deletePushSubscriptionHandler(w, r)
	}
}

func getPushSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	endpoint := query.Get("endpoint")
	if endpoint == "" {
		http.Error(w, InvalidEndPointQuery, http.StatusInternalServerError)
		return
	}
}

type PostPushSubscriptionBody struct {
	Endpoint       string     `validate:"required"`
	P256hd         string     `validate:"required"`
	Auth           string     `validate:"required"`
	ExpirationTime *time.Time `validate:"required"`
}

func postPushSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
}

func deletePushSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
}

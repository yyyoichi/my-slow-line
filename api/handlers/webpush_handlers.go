package handlers

import (
	"encoding/json"
	"fmt"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"net/http"
	"os"
	"strconv"
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
	ExpirationTime *time.Time `validate:""`
	UserAgent      string     `validate:"required"`
}

func postPushSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &PostPushSubscriptionBody{}
	if err := utils.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read context
	userId, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set subscribe
	webp := services.NewRepositoryServices().GetWebpush()
	if err = webp.Subscription(userId, b.Endpoint, b.P256hd, b.Auth, b.UserAgent, b.ExpirationTime); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send notifier
	if err = webp.SendNotification(userId, "Thanks!"); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deletePushSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
}

package handlers

import (
	"encoding/json"
	"fmt"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"himakiwa/services/webpush"
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

type WebPushSubscriptionHandlers struct {
	services.UseRepositoryServices
	webpush.UserWebpushServices
}

func NewWebPushSubscriptionHandlers(useRepository services.UseRepositoryServices, useWebpush webpush.UserWebpushServices) func(http.ResponseWriter, *http.Request) {
	wsh := &WebPushSubscriptionHandlers{useRepository, useWebpush}
	return wsh.PushSubscriptionHandler
}

func (wsh *WebPushSubscriptionHandlers) PushSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		wsh.getPushSubscriptionHandler(w, r)
	case "POST":
		wsh.postPushSubscriptionHandler(w, r)
	case "DELETE":
		wsh.deletePushSubscriptionHandler(w, r)
	}
}

func (wsh *WebPushSubscriptionHandlers) getPushSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
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

func (wsh *WebPushSubscriptionHandlers) postPushSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &PostPushSubscriptionBody{}
	if err := utils.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set subscribe
	if err := wsh.UseRepositoryServices(userID).UserServices.AddWebpushSubscription(b.Endpoint, b.P256hd, b.Auth, b.UserAgent, b.ExpirationTime); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send webpush
	if err := wsh.UserWebpushServices(b.Endpoint, b.Auth, b.P256hd).SendPlaneMessage("Thanks!"); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (wsh *WebPushSubscriptionHandlers) deletePushSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
}

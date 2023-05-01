package handlers

import (
	"encoding/json"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrUnExpcetedMethod = "invalid method"
	ErrSystem           = "error in system"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getMeHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	utils.DeleteJWTCookie(w)
	w.WriteHeader(http.StatusOK)
}

type GetMeResp struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	LoginAt        time.Time `json:"loginAt"`
	CreateAt       time.Time `json:"createAt"`
	UpdateAt       time.Time `json:"updateAt"`
	TwoVerificated bool      `json:"towVerificatedAt"`
	Deleted        bool      `json:"deleted"`
}

func getMeHandler(w http.ResponseWriter, r *http.Request) {
	// read context
	userId, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ur := services.NewRepositoryServices().GetUser()
	tu, err := ur.Query(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GetMeResp{
		tu.Id,
		tu.Name,
		tu.Email,
		tu.LoginAt,
		tu.CreateAt,
		tu.UpdateAt,
		tu.TwoVerificated,
		tu.Deleted,
	}

	// set jwt-token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

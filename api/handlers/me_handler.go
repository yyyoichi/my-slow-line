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

type MeHandlers struct {
	services.UseRepositoryServices
}

func NewMeHandlers(use services.UseRepositoryServices) *MeHandlers {
	return &MeHandlers{use}
}

func (mh *MeHandlers) MeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		mh.getMeHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

func (*MeHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	utils.DeleteJWTCookie(w)
	w.WriteHeader(http.StatusOK)
}

type GetMeResp struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	LoginAt        time.Time `json:"loginAt"`
	CreateAt       time.Time `json:"createAt"`
	UpdateAt       time.Time `json:"updateAt"`
	TwoVerificated bool      `json:"towVerificatedAt"`
	Deleted        bool      `json:"deleted"`
}

func (mh *MeHandlers) getMeHandler(w http.ResponseWriter, r *http.Request) {
	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userServices := mh.UseRepositoryServices(userID).UserServices
	user, err := userServices.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GetMeResp{
		user.ID,
		user.Name,
		user.Email,
		user.LoginAt.Time,
		user.CreateAt,
		user.UpdateAt,
		user.TwoVerificated,
		user.Deleted,
	}

	// set jwt-token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

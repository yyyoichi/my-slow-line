package handlers

import (
	"encoding/json"
	"errors"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	ErrNotExistRecruit = errors.New("dose not exist recruit")
)

type PublicRecruitHandlers struct {
	services.UseRepositoryServices
}

func NewPublicRecruitHandlers(use services.UseRepositoryServices) func(w http.ResponseWriter, r *http.Request) {
	prh := &PublicRecruitHandlers{use}
	return prh.PublicRecruitmentHandler
}

func (prh *PublicRecruitHandlers) PublicRecruitmentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		prh.getRecruitmentHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

type GetPublicRecruitmentResp struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Uuid    string `json:"uuid"`
}

func (prh *PublicRecruitHandlers) getRecruitmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recruitmentUUID := vars["recruitmentUUID"]
	if recruitmentUUID == "" {
		http.Error(w, ErrNotExistRecruit.Error(), http.StatusInternalServerError)
		return
	}

	// get reqruit by uuid
	user, err := prh.UseRepositoryServices(0).UserServices.GetUserByRecruitUUID(recruitmentUUID)
	if err != nil {
		http.Error(w, ErrSystem, http.StatusInternalServerError)
		return
	}

	// map
	resp := GetPublicRecruitmentResp{
		user.Name,
		user.Message,
		user.UUID,
	}

	// resp
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type RecruitHandlers struct {
	services.UseRepositoryServices
}

func NewRecruitHandlers(use services.UseRepositoryServices) func(w http.ResponseWriter, r *http.Request) {
	rh := &RecruitHandlers{use}
	return rh.RecruitmentHandler
}

func (rh *RecruitHandlers) RecruitmentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rh.GetRecruitmentsHandler(w, r)
	case "POST":
		rh.PostRecruitmentsHandler(w, r)
	case "PUT":
		rh.PutRecruitmentsHandler(w, r)
	case "DELETE":
		rh.DeleteRecruitmentsHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

type GetRecruitmentResp struct {
	Id       int       `json:"name"`
	UserId   int       `json:"userID"`
	Uuid     string    `json:"uuid"`
	Message  string    `json:"message"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
	Deleted  bool      `json:"deleted"`
}

func (rh *RecruitHandlers) GetRecruitmentsHandler(w http.ResponseWriter, r *http.Request) {
	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get service
	userServices := rh.UseRepositoryServices(userID).UserServices

	// get data
	recruits, err := userServices.GetRecruitments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// map
	recruitsments := []GetRecruitmentResp{}
	for _, r := range recruits {
		recruitsments = append(recruitsments, GetRecruitmentResp{
			Id: r.ID, UserId: r.UserID, Uuid: r.UUID, Message: r.Message, CreateAt: r.CreateAt, UpdateAt: r.UpdateAt, Deleted: r.Deleted,
		})
	}

	//resp
	err = json.NewEncoder(w).Encode(recruitsments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type PostRecruitmentBody struct {
	Message string `validate:"required"`
}

// create recruitments
func (rh *RecruitHandlers) PostRecruitmentsHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &PostRecruitmentBody{}
	if err := utils.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get service
	userServices := rh.UseRepositoryServices(userID).UserServices

	// create
	if err = userServices.CreateRecruitment(b.Message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type PutRecruitmentBody struct {
	Uuid    string `validate:"required"`
	Message string `validate:"required"`
	Deleted bool   `validate:"required"`
}

// udate recruitments
func (rh *RecruitHandlers) PutRecruitmentsHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &PutRecruitmentBody{}
	if err := utils.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get service
	userServices := rh.UseRepositoryServices(userID).UserServices

	// update
	if err = userServices.UpdateRecruitment(b.Uuid, b.Message, b.Deleted); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// delete recruitments
func (rh *RecruitHandlers) DeleteRecruitmentsHandler(w http.ResponseWriter, r *http.Request) {
}

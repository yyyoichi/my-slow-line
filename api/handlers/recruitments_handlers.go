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

func PublicRecruitmentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getRecruitmentHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

var (
	ErrNotExistRecruit = errors.New("dose not exist recruit")
)

type GetPublicRecruitmentResp struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Uuid    string `json:"uuid"`
}

func getRecruitmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recruitmentUUID := vars["recruitmentUUID"]
	if recruitmentUUID == "" {
		http.Error(w, ErrNotExistRecruit.Error(), http.StatusInternalServerError)
		return
	}

	users := services.NewRepositoryServices().GetUser()

	user, err := users.QueryByRecruitUuid(recruitmentUUID)
	if err != nil {
		http.Error(w, ErrSystem, http.StatusInternalServerError)
		return
	}
	resp := GetPublicRecruitmentResp{
		user.Name,
		user.Message,
		user.Uuid,
	}
	json.NewEncoder(w).Encode(resp)
}

func RecruitmentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getRecruitmentsHandler(w, r)
	case "POST":
		postRecruitmentsHandler(w, r)
	case "PUT":
		putRecruitmentsHandler(w, r)
	case "DELETE":
		deleteRecruitmentsHandler(w, r)
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

func getRecruitmentsHandler(w http.ResponseWriter, r *http.Request) {
	// read context
	userId, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users := services.NewRepositoryServices().GetUser()
	recruit := users.GetFriendRecruitService(userId)

	// get data
	recruits, err := recruit.Query()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// map
	recruitsments := []GetRecruitmentResp{}
	for _, r := range recruits {
		recruitsments = append(recruitsments, GetRecruitmentResp{
			Id: r.Id, UserId: r.UserId, Uuid: r.Uuid, Message: r.Message, CreateAt: r.CreateAt, UpdateAt: r.UpdateAt, Deleted: r.Deleted,
		})
	}

	//resp
	err = json.NewEncoder(w).Encode(recruitsments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// create recruitments
func postRecruitmentsHandler(w http.ResponseWriter, r *http.Request) {
}

// udate recruitments
func putRecruitmentsHandler(w http.ResponseWriter, r *http.Request) {
}

// delete recruitments
func deleteRecruitmentsHandler(w http.ResponseWriter, r *http.Request) {
}

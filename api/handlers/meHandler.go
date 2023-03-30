package handlers

import (
	"encoding/json"
	"himakiwa/database"
	"himakiwa/utils"
	"net/http"
	"strconv"
	"time"
)

type MeJSON struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	LoginAt  time.Time `json:"loginAt"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
	Deleted  bool      `json:"deleted"`
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	// read context
	userId, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get from db
	du, err := database.QueryUser(nil, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// map and response
	u := &MeJSON{du.Id, du.Name, du.Email, du.LoginAt, du.CreateAt, du.UpdateAt, du.Deleted}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

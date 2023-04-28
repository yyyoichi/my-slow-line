package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"himakiwa/auth"
	"himakiwa/handlers/decode"
	"himakiwa/services/database"
	"himakiwa/utils"
	"net/http"
	"os"
)

var (
	ErrInvalidCode = errors.New("invalid twe step verification code")
)

type VerificationCodeBody struct {
	Code string `validate:"required,len=6"`
	Id   int    `validate:"required"`
}

// verficate code return jwt in set-cookie
func VerificateTwoStepCodeHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &VerificationCodeBody{}
	if err := decode.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId := b.Id

	// connect db
	db := database.DB

	// get user
	du, err := database.QueryUser(db, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// compare code
	result := auth.VerificateSixNumber(b.Code, du.TwoStepVerificationCode, du.LoginAt)
	if !result {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrInvalidCode)
		return
	}

	// set-cookie jwt-token
	jt := auth.NewJwt(os.Getenv("JWT_SECRET"))
	token, err := jt.Generate(fmt.Sprintf("%d", userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SetJWTCookie(w, token)
	w.WriteHeader(http.StatusOK)
}

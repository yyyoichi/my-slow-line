package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrInvalidBodyProperty = errors.New("invalid body property")
	ErrInvalidBody         = errors.New("invalid body")
)

func DecodeBody(r *http.Request, expStruct interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(expStruct); err != nil {
		return ErrInvalidBodyProperty
	}

	vld := newValidation()
	if err := vld.Struct(expStruct); err != nil {
		return ErrInvalidBody
	}
	return nil
}

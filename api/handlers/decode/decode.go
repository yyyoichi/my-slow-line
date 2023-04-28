package decode

import (
	"encoding/json"
	"net/http"
)

func DecodeBody(r *http.Request, expStruct interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(expStruct); err != nil {
		return errInvalidBodyProperty
	}

	vld := newValidation()
	if err := vld.Struct(expStruct); err != nil {
		return errInvalidBody
	}
	return nil
}

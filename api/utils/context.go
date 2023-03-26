package utils

import (
	"context"
	"net/http"
)

type UserKey string

const (
	userId UserKey = "claims"
)

func WithUserContext(r *http.Request, id string) *http.Request {
	cxt := context.WithValue(r.Context(), userId, id)
	return r.WithContext(cxt)
}
func ReadUserContext(r *http.Request) string {
	u, ok := r.Context().Value(userId).(string)
	if !ok {
		return ""
	}
	return u
}

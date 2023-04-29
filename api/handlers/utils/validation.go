package utils

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	errInvalidBodyProperty = errors.New("invalid body property")
	errInvalidBody         = errors.New("invalid body")
)

func newValidation() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("alphanumary", customAlphanumary)
	return validate
}

func customAlphanumary(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	rg := regexp.MustCompile(`([0-9].*[a-zA-Z]|[a-zA-Z].*[0-9])`)
	return rg.MatchString(s)
}

package auth

import (
	"testing"
)

func TestSafe(t *testing.T) {
	tests := []string{
		"abcd",
		"testpass",
		"_gg",
	}
	for _, tt := range tests {
		hashedPassword, err := HashPassword(tt)
		if err != nil {
			t.Error(err)
		}
		result, err := ComparePasswordAndHash(tt, hashedPassword)
		if err != nil {
			t.Errorf("expected '%s' but got not equar", tt)
		}
		if !result {
			t.Errorf("does not match password '%s'", tt)
		}
	}
}

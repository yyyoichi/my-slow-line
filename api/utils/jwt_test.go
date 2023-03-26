package utils

import (
	"errors"
	"testing"
	"time"
)

func newTestJWT(secret string, start time.Time) *JwtToken {
	return &JwtToken{secret, start}
}
func TestJwt(t *testing.T) {
	tests := []struct {
		jwt    *JwtToken
		genid  string
		prssc  string
		experr error
		expid  string
	}{
		{ // 正常系
			jwt:    newTestJWT("abc", time.Now()),
			genid:  "testid",
			prssc:  "abc",
			experr: nil,
			expid:  "testid",
		},
		{
			//異常系（期限切れ）
			jwt:    newTestJWT("abc", time.Now().Add(time.Duration(-8)*time.Hour*24)),
			genid:  "testid",
			prssc:  "abc",
			experr: errors.New("invalid"),
			expid:  "",
		},
		{
			//異常系（ソルト）
			jwt:    newTestJWT("abc", time.Now()),
			genid:  "testid",
			prssc:  "ffff",
			experr: errors.New("invalid"),
			expid:  "",
		},
	}
	for _, tt := range tests {
		token, err := tt.jwt.Generate(tt.genid)
		if err != nil {
			t.Errorf("cannot generate token")
		}
		jwt := NewJwt(tt.prssc)
		rc, err := jwt.ParseToken(token)

		if err == nil && tt.experr != nil {
			t.Errorf("expected err='%s' but actual got='nil'", tt.experr)
		}
		if err != nil && tt.experr == nil {
			t.Errorf("unexpected err. expected='nil' but got=%s", err)
		}

		if rc == nil && tt.expid == "" {
			continue
		}

		if rc.ID != tt.expid {
			t.Errorf("expected id id '%s' but got=%s", tt.expid, rc.ID)
		}
	}
}

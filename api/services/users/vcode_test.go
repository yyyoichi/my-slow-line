package users

import (
	"testing"
	"time"
)

func TestTwoVerification(t *testing.T) {
	input := generateRandomSixNumber()
	test := []struct {
		code string
		add  time.Duration
		exp  bool
	}{

		{ // normal
			code: input,
			add:  -time.Minute * 9,
			exp:  true,
		},
		{ // normal
			code: input,
			add:  -time.Second * 599,
			exp:  true,
		},
		{ // abnormal
			code: input,
			add:  -time.Second * 600,
			exp:  false,
		},
		{ // abnormal
			code: generateRandomSixNumber(),
			add:  -time.Minute * 9,
			exp:  false,
		},
		{ // abnormal
			code: generateRandomSixNumber(),
			add:  -time.Second * 601,
			exp:  false,
		},
		{ // abnormal
			code: generateRandomSixNumber(),
			add:  time.Second * 1,
			exp:  false,
		},
	}
	for i, tt := range test {
		loginAt := time.Now().Add(tt.add)
		result := verificateSixNumber(input, tt.code, loginAt)
		if result != tt.exp {
			t.Errorf("%d: expected '%v' but got='%v' code is '%s' loginAt is before '%v'", i, tt.exp, result, tt.code, tt.add)
		}
	}
}

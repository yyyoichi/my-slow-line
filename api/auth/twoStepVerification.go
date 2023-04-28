package auth

import (
	"math/rand"
	"time"
)

func GenerateRandomSixNumber() string {
	rand.Seed(time.Now().UnixNano())
	numbers := []byte("0123456789")
	result := make([]byte, 6)
	for i := 0; i < 6; i++ {
		result[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(result)
}

// verification is success when a [in] code from user form match [code] from db and [loginAt] time is within 10min.
func VerificateSixNumber(in, code string, loginAt time.Time) bool {
	if in != code {
		return false
	}
	now := time.Now()
	max := loginAt.Add(time.Minute * 10)
	return now.Before(max)
}

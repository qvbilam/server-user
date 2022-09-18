package utils

import (
	"math/rand"
	"strings"
	"time"
)

const (
	randomDefaultAlphabet       = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
	randomDeFaultNumberAlphabet = "1234567890"
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func Random(length int, alphabet string) string {
	if length <= 0 {
		return ""
	}
	if alphabet == "" {
		alphabet = randomDefaultAlphabet
	}

	source := rand.NewSource(time.Now().UnixNano())
	s := strings.Builder{}
	s.Grow(length)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := length-1, source.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = source.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(alphabet) {
			s.WriteByte(alphabet[idx])
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return s.String()
}

func RandomCharacter(length int) string {
	return Random(length, randomDefaultAlphabet)
}

func RandomNumber(length int) string {
	return Random(length, randomDeFaultNumberAlphabet)
}

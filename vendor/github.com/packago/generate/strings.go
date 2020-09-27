package generate

import (
	"math/rand"
	"time"
)

const (
	lettersNumbersBytes   = "abcdefghkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	lowercaseNumbersBytes = "abcdefghkmnopqrstuvwxyz23456789"
	letterIdxBits         = 6
	letterIdxMask         = 1<<letterIdxBits - 1
	letterIdxMax          = 63 / letterIdxBits
)

var randUnixNanoSrc = rand.NewSource(time.Now().UnixNano())

func RandomLettersNumbers(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, randUnixNanoSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randUnixNanoSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(lettersNumbersBytes) {
			b[i] = lettersNumbersBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func RandomLowercaseNumbers(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, randUnixNanoSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randUnixNanoSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(lowercaseNumbersBytes) {
			b[i] = lowercaseNumbersBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

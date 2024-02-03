package utils

import (
	"math/rand"
	"time"
	"unsafe"
)

type Charset string

var (
	CHARSET_NUM         Charset = "0123456789"
	CHARSET_ALPHA_LOWER Charset = "abcdefghijklmnopqrstuvwxyz"
	CHARSET_ALPHA_UPPER Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CHARSET_ALPHA       Charset = CHARSET_ALPHA_LOWER + CHARSET_ALPHA_UPPER
	CHARSET_ALPHA_NUM   Charset = CHARSET_ALPHA + CHARSET_NUM
	CHARSET_URL_SAFE    Charset = "_-" + CHARSET_ALPHA_NUM
)

func NewRandomStringGenerator(charset Charset, length int) func() string {

	var prngSrc = rand.NewSource(time.Now().UnixNano())

	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	return func() string {

		b := make([]byte, length)
		// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
		for i, cache, remain := length-1, prngSrc.Int63(), letterIdxMax; i >= 0; {
			if remain == 0 {
				cache, remain = prngSrc.Int63(), letterIdxMax
			}
			if idx := int(cache & letterIdxMask); idx < len(charset) {
				b[i] = charset[idx]
				i--
			}
			cache >>= letterIdxBits
			remain--
		}

		return *(*string)(unsafe.Pointer(&b))
	}
}

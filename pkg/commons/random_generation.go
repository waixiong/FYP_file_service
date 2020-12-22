package commons

import (
	mathRand "math/rand"
	"time"
	"unsafe"
)

// ---------------- Random String Genarator ---------------- //
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const digits = "1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

	digitIdxBits = 4                   // 6 bits to represent a letter index
	digitIdxMask = 1<<digitIdxBits - 1 // All 1-bits, as many as letterIdxBits
	digitIdxMax  = 10 / digitIdxBits   // # of letter indices fitting in 63 bits
)

var src = mathRand.NewSource(time.Now().UnixNano())

// source : https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func randDigitBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), digitIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), digitIdxMax
		}
		if idx := int(cache & digitIdxMask); idx < len(digits) {
			b[i] = digits[idx]
			i--
		}
		cache >>= digitIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

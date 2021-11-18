package methods

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"time"
)

// Sign is a method for HMAC-SHA1 signer
func Sign(signingKey, message string) string {
	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(message))
	signatureBytes := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signatureBytes)
}

// RandomString is used for random string generation on basis of length
func RandomString(l int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt())
	}
	return string(bytes)
}

// generating random integer
func randInt() int {
	rdInt := rand.Intn(51)
	if rdInt <= 25 {
		rdInt = rdInt + 65
	} else {
		rdInt = rdInt + 71
	}
	return rdInt
}

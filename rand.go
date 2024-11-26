package misc

import (
	"crypto/rand"
)

func RandomString(length int, dic string) string {
	b := make([]byte, length)
	rand.Read(b)
	dicLen := byte(len(dic))
	for i := 0; i < length; {
		b[i] = dic[b[i]%dicLen]
		i++
	}
	return string(b)
}

// RandomStringAlphabet
func RandomStringA(length int) string {
	return RandomString(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// RandomStringNumeric
func RandomStringN(length int) string {
	return RandomString(length, "0123456789")
}

// RandomStringAlphabetNumeric
func RandomStringAN(length int) string {
	return RandomString(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

// RandomStringAlphabetNumericSymbol
func RandomStringANS(length int) string {
	return RandomString(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}\\|;:'\",./<>?")
}

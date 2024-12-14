package misc

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
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

func RandomUInt64() uint64 {
	b := make([]byte, 8)
	rand.Read(b)
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func RandomInt64() int64 {
	b := make([]byte, 8)
	rand.Read(b)
	var v int64
	binary.Read(bytes.NewReader(b), binary.LittleEndian, &v)
	return v
}

func RandomUInt32() uint32 {
	b := make([]byte, 4)
	rand.Read(b)
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func RandomInt32() int32 {
	b := make([]byte, 4)
	rand.Read(b)
	var v int32
	binary.Read(bytes.NewReader(b), binary.LittleEndian, &v)
	return v
}

func RandomInt() int {
	return int(RandomInt64())
}

func RandomIntRange(min, max int) int {
	r := RandomInt()
	if r < 0 {
		r = -r
	}
	return min + r%(max-min)
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

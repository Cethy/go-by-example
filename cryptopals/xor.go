package main

import (
	"encoding/hex"
)

// https://cryptopals.com/sets/1/challenges/2

func xor(s string, x string) (string, error) {
	sDecoded, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}
	xDecoded, err := hex.DecodeString(x)
	if err != nil {
		return "", err
	}

	r := make([]byte, len(sDecoded))
	for i := 0; i < len(sDecoded); i++ {
		r[i] = sDecoded[i] ^ xDecoded[i]
	}

	return hex.EncodeToString(r), nil
}

// https://cryptopals.com/sets/1/challenges/3

func singleByteXor(s string, x byte) (string, error) {
	r := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		r[i] = s[i] ^ x
	}

	return string(r), nil
}

// https://cryptopals.com/sets/1/challenges/5

func repeatingBytesXor(s string, x []byte) (string, error) {
	r := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		r[i] = s[i] ^ x[i%len(x)]
	}

	return string(r), nil
}

package main

import (
	"encoding/base64"
	"encoding/hex"
)

// https://cryptopals.com/sets/1/challenges/1

func hexToBase64(s string) (string, error) {
	sDecoded, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sDecoded), nil
}

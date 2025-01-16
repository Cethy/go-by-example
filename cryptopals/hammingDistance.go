package main

import (
	"errors"
	"strconv"
	"strings"
)

func hammingDistance(s1, s2 string) (int, error) {
	if len(s1) != len(s2) {
		return -1, errors.New("input strings must have same length")
	}

	distance := 0
	for i := 0; i < len(s1); i++ {
		r := s1[i] ^ s2[i]
		binaryStr := strconv.FormatUint(uint64(r), 2)
		distance += strings.Count(binaryStr, "1")
	}

	return distance, nil
}

package main

import (
	"math"
	"strings"
)

func scorePlaintext(plaintext string) int {
	freq := make(map[rune]int)
	for _, char := range plaintext {
		if char >= ' ' {
			freq[char]++
		}
		if char >= 'a' && char <= 'z' {
			freq[char]++
		}
	}

	refFreq := " etaoinshrdlcumwfgypbvkjxqz"
	score := 0
	for _, char := range refFreq {
		score += int(math.Abs(float64(freq[char] - strings.Count(refFreq, string(char)))))
	}
	return score
}

package main

import "sort"

/**
 * Only returns best candidate
 */
func breakSingleByteXor(source string) (string, string, int) {
	var scores []struct {
		key     string
		decoded string
		score   int
	}

	for i := 0; i < 129; i++ {
		r, err := singleByteXor(source, byte(i))
		if err != nil {
			panic(err)
		}

		score := scorePlaintext(r)
		scores = append(scores, struct {
			key     string
			decoded string
			score   int
		}{key: string(byte(i)), score: score, decoded: r})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	return scores[0].key, scores[0].decoded, scores[0].score
}

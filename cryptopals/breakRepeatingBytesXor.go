package main

func breakRepeatingBytesXor(source string, keySize int) (string, string) {
	// make a block for each char of the key
	blocks := make([]string, keySize)
	for i := 0; i < len(source); i += keySize {
		end := i + keySize
		if end > len(source) {
			end = len(source)
		}

		for j := 0; j < keySize; j++ {
			if len(source) <= i+j {
				break
			}
			blocks[j] += string(source[i:end][j])
		}
	}

	decodedBlocks := make([]string, keySize)

	// using breaker
	key := ""
	for i, block := range blocks {
		keyPart, decoded, _ := breakSingleByteXor(block)
		//println(foundKey, decoded, score)
		decodedBlocks[i] = decoded
		key += keyPart
	}

	// rebundling blocks
	repacked := ""
	for i := 0; i < len(decodedBlocks[0]); i++ {
		for _, block := range decodedBlocks {
			if len(block) > i {
				repacked += string(block[i])
			}
		}
	}

	return key, repacked
}

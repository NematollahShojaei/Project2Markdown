// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package utils

import "unicode"

// EstimateTokens provides a fast, heuristic-based token count without heavy ML libraries.
// It counts words, symbols, and punctuation which closely mimics BPE tokenizers for code.
func EstimateTokens(s string) int {
	count := 0
	inWord := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			inWord = false
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			if !inWord {
				count++
				inWord = true
			}
		} else {
			// Punctuation and symbols count as individual tokens
			count++
			inWord = false
		}
	}
	// Add 1 token for the line break/structure
	return count + 1
}

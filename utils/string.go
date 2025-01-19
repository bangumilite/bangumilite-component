package utils

import (
	"github.com/agnivade/levenshtein"
	"golang.org/x/text/unicode/norm"
)

func IsMatch(a, b string, s float64) bool {
	// normalize strings to handle differences in encoding or forms
	a = norm.NFKC.String(a)
	b = norm.NFKC.String(b)

	if a == b {
		return true
	}

	distance := levenshtein.ComputeDistance(a, b)

	// calculate similarity based on rune counts
	runeLenA := len([]rune(a))
	runeLenB := len([]rune(b))
	longerLength := runeLenA
	if runeLenB > longerLength {
		longerLength = runeLenB
	}

	// avoid division by zero if both strings are empty after normalization
	if longerLength == 0 {
		return true
	}

	// calculate and return similarity percentage
	similarity := (1.0 - float64(distance)/float64(longerLength)) * 100
	if similarity < s {
		return false
	}

	return true
}

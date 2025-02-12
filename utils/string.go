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

	runeLenA := len([]rune(a))
	runeLenB := len([]rune(b))
	longerLength := runeLenA
	if runeLenB > longerLength {
		longerLength = runeLenB
	}

	if longerLength == 0 {
		return true
	}

	distance := levenshtein.ComputeDistance(a, b)
	similarity := (1.0 - float64(distance)/float64(longerLength)) * 100
	if similarity < s {
		return false
	}

	return true
}

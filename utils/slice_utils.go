package utils

// RemoveDuplicates removes duplicates from a slice of any comparable type.
func RemoveDuplicates[T comparable](elems []T) []T {
	seen := make(map[T]bool)
	var res []T

	for _, e := range elems {
		if !seen[e] {
			seen[e] = true
			res = append(res, e)
		}
	}

	return res
}

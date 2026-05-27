package edlib

// QgramDistance compute the q-gram similarity between two strings
// Takes two strings as parameters, a split length which defines the k-gram shingle length
func QgramDistance(str1, str2 string, splitLength int) int { _ = "STUB: not implemented"; return 0 }

// QgramDistanceCustomNgram compute the q-gram similarity between two custom set of individuals
// Takes two n-gram map as parameters
func QgramDistanceCustomNgram(splittedStr1, splittedStr2 map[string]int) int {
	_ = "STUB: not implemented"
	return 0
}

// QgramSimilarity compute a similarity index (between 0 and 1) between two strings from a Qgram distance
// Takes two strings as parameters, a split length which defines the k-gram shingle length
func QgramSimilarity(str1, str2 string, splitLength int) float32 {
	_ = "STUB: not implemented"
	return 0
}

package edlib

// JaccardSimilarity compute the jaccard similarity coeffecient between two strings
// Takes two strings as parameters, a split length which define the k-gram single length
// (if zero split string on whitespaces) and return an index.
func JaccardSimilarity(str1, str2 string, splitLength int) float32 {
	_ = "STUB: not implemented"
	return 0
}

// Split string before rune conversion for jaccard calculation
// If splitLength == 0 then split on whitespaces
// Else use shingle algorithm

// Conversion of splitted string into rune array

// Create union keywords slice between input strings

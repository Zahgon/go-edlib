package edlib

// CosineSimilarity use cosine algorithm to return a similarity index between string vectors
// Takes two strings as parameters, a split length which define the k-gram single length
// (if zero split string on whitespaces) and return an index.
func CosineSimilarity(str1, str2 string, splitLength int) float32 {
	_ = "STUB: not implemented"
	return 0
}

// Split string before rune conversion for cosine calculation
// If splitLength == 0 then split on whitespaces
// Else use shingle algorithm

// Conversion of plitted string into rune array

// Create union keywords slice between input strings

// Compute cosine algorithm

// Compute union between two string slices, convert result to rune matrix and return it
func union(a, b []string) [][]rune { _ = "STUB: not implemented"; return nil }

// Convert a to rune matrix (with x -> words and y -> characters)

// Find takes a rune slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1.
func find(slice [][]rune, val []rune) int { _ = "STUB: not implemented"; return 0 }

// Return the elements sum from int slice
func sum(arr []int) int { _ = "STUB: not implemented"; return 0 }

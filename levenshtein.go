package edlib

// LevenshteinDistance calculate the distance between two string
// This algorithm allow insertions, deletions and substitutions to change one string to the second
// Compatible with non-ASCII characters
func LevenshteinDistance(str1, str2 string) int {
	_ = "STUB: not implemented"
	// Convert string parameters to rune arrays to be compatible with non-ASCII
	return 0
}

// Get and store length of these strings

// insert
// delete
// substitution

// OSADamerauLevenshteinDistance calculate the distance between two string
// Optimal string alignment distance variant that use extension of the Wagner-Fisher dynamic programming algorithm
// Doesn't allow multiple transformations on a same substring
// Allowing insertions, deletions, substitutions and transpositions to change one string to the second
// Compatible with non-ASCII characters
func OSADamerauLevenshteinDistance(str1, str2 string) int {
	_ = "STUB: not implemented"
	// Convert string parameters to rune arrays to be compatible with non-ASCII
	return 0
}

// Get and store length of these strings

// 2D Array

// insertion, deletion, substitution

// translation

// DamerauLevenshteinDistance calculate the distance between two string
// This algorithm computes the true Damerau–Levenshtein distance with adjacent transpositions
// Allowing insertions, deletions, substitutions and transpositions to change one string to the second
// Compatible with non-ASCII characters
func DamerauLevenshteinDistance(str1, str2 string) int {
	_ = "STUB: not implemented"
	// Convert string parameters to rune arrays to be compatible with non-ASCII
	return 0
}

// Get and store length of these strings

// Create alphabet based on input strings

// 2D Array for distance matrix : matrix[0..str1.length+2][0..s2.length+2]

// Maximum possible distance

// Initialize matrix

// Process edit distance

// Addition
// Deletion

// Substitution
// Transposition

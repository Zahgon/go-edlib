package edlib

// JaroSimilarity return a similarity index (between 0 and 1)
// It use Jaro distance algorithm and allow only transposition operation
func JaroSimilarity(str1, str2 string) float32 {
	_ = "STUB: not implemented"
	// Convert string parameters to rune arrays to be compatible with non-ASCII
	return 0
}

// Get and store length of these strings

// Maximum matching distance allowed

// Correspondence tables (1 for matching and 0 if it's not the case)

// Check for matching characters in both strings

// Check for possible translations

// JaroWinklerSimilarity return a similarity index (between 0 and 1)
// Use Jaro similarity and after look for a common prefix (length <= 4)
func JaroWinklerSimilarity(str1, str2 string) float32 {
	_ = "STUB: not implemented"
	// Get Jaro similarity index between str1 and str2
	return 0
}

// Convert string parameters to rune arrays to be compatible with non-ASCII

// Get and store length of these strings

// Find length of the common prefix

// Normalized prefix count with Winkler's constraint
// (prefix length must be inferior or equal to 4)

// Return calculated Jaro-Winkler similarity index

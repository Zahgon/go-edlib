package edlib

// Algorithm is an Integer type used to identify edit distance algorithms
type Algorithm uint8

// Algorithm identifiers
const (
	Levenshtein Algorithm = iota
	DamerauLevenshtein
	OSADamerauLevenshtein
	Lcs
	Hamming
	Jaro
	JaroWinkler
	Cosine
	Jaccard
	SorensenDice
	Qgram
)

// StringsSimilarity return a similarity index [0..1] between two strings based on given edit distance algorithm in parameter.
// Use defined Algorithm type.
// Through this function, Cosine and Jaccard algorithms are used with Shingle split method with a length of 2.
func StringsSimilarity(str1 string, str2 string, algo Algorithm) (float32, error) {
	_ = "STUB: not implemented"
	return 0, nil
}

// Return matching index E [0..1] from two strings and an edit distance
func matchingIndex(str1 string, str2 string, distance int) float32 {
	_ = "STUB: not implemented"
	// Convert strings to rune slices
	return 0
}

// Compare rune arrays length and make a matching percentage between them

// FuzzySearch realize an approximate search on a string list and return the closest one compared
// to the string input
func FuzzySearch(str string, strList []string, algo Algorithm) (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

// FuzzySearchThreshold realize an approximate search on a string list and return the closest one compared
// to the string input. Takes a similarity threshold in parameter.
func FuzzySearchThreshold(str string, strList []string, minSim float32, algo Algorithm) (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

// FuzzySearchSet realize an approximate search on a string list and return a set composed with x strings compared
// to the string input sorted by similarity with the base string.
// Takes the a quantity parameter to define the number of output strings desired
// (For example 3 in the case of the Google Keyboard word suggestion).
func FuzzySearchSet(str string, strList []string, quantity int, algo Algorithm) ([]string, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// FuzzySearchSetThreshold realize an approximate search on a string list and return a set composed with x strings compared
// to the string input sorted by similarity with the base string. Take a similarity threshold in parameter.
// Takes the a quantity parameter to define the number of output strings desired
// (For example 3 in the case of the Google Keyboard word suggestion).
// Takes also a threshold parameter for similarity with base string.
func FuzzySearchSetThreshold(str string, strList []string, quantity int, minSim float32, algo Algorithm) ([]string, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

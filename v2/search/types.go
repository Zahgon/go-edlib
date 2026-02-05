package search

import "github.com/hbollon/go-edlib/v2/algorithms"

// Result represents a single search result
type Result struct {
	// Value is the matched string
	Value string

	// Score is the similarity score [0.0 to 1.0]
	Score float64

	// Distance is the edit distance from the query
	Distance int

	// Index is the original index in the input slice
	Index int

	// Metadata contains additional information about the result
	Metadata ResultMetadata
}

// ResultMetadata contains metadata about a search result
type ResultMetadata struct {
	// Rank is the ranking position (1-based)
	Rank int

	// Algorithm used for the search
	Algorithm algorithms.Algorithm

	// AboveThreshold indicates if the result met the threshold
	AboveThreshold bool
}

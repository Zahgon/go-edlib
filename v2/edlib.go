package edlib

// Package edlib provides string comparison and edit distance algorithms
// with a modern, fluent API design.
//
// Basic usage:
//
//	// Quick one-liners
//	distance, _ := edlib.Quick().Levenshtein("hello", "helo")
//	similarity, _ := edlib.Quick().JaroWinkler("test", "testing")
//
//	// Fluent API with configuration
//	result, _ := edlib.New().
//		Using(algorithms.Levenshtein).
//		WithThreshold(0.8).
//		WithCaseSensitive(false).
//		Compare("Hello", "hello")
//
//	// Fuzzy search
//	results, _ := edlib.New().
//		Search().
//		WithAlgorithm(algorithms.JaroWinkler).
//		WithThreshold(0.7).
//		FindAll("test", []string{"tester", "testing", "best"})
//
// Supported algorithms:
//   - Edit Distance: Levenshtein, Damerau-Levenshtein, Hamming
//   - Similarity: Jaro, Jaro-Winkler, Cosine, Jaccard, Sorensen-Dice
//   - Sequence: LCS (Longest Common Subsequence)
//   - Phonetic: Soundex, Metaphone, Double Metaphone
//
// Features:
//   - Zero dependencies (except OpenTelemetry for optional metrics)
//   - Unicode support
//   - Memory pooling for performance
//   - Parallel processing
//   - Context support for cancellation
//   - Detailed execution metrics
//   - Type-safe fluent API

import (
	// Import algorithm implementations to trigger auto-registration
	_ "github.com/hbollon/go-edlib/v2/internal/impl/edit_distance/levenshtein"
	_ "github.com/hbollon/go-edlib/v2/internal/impl/similarity/cosine"
	_ "github.com/hbollon/go-edlib/v2/internal/impl/similarity/jaro"
)

// Version is the current version of the library
const Version = "2.0.0"

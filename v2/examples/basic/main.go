package main

import (
	"fmt"

	"github.com/hbollon/go-edlib/v2"
	"github.com/hbollon/go-edlib/v2/algorithms"
)

func main() {
	fmt.Println("=== go-edlib v2.0 Examples ===\n")

	// Example 1: Quick API for simple operations
	quickExamples()

	// Example 2: Fluent API with configuration
	fluentExamples()

	// Example 3: Fuzzy search
	searchExamples()

	// Example 4: Batch processing
	batchExamples()

	// Example 5: Advanced features
	advancedExamples()
}

func quickExamples() {
	fmt.Println("1. Quick API Examples:")
	fmt.Println("----------------------")

	// Levenshtein distance
	dist, _ := edlib.Quick().Levenshtein("kitten", "sitting")
	fmt.Printf("Levenshtein('kitten', 'sitting') = %d\n", dist)

	// Jaro-Winkler similarity
	sim, _ := edlib.Quick().JaroWinkler("martha", "marhta")
	fmt.Printf("JaroWinkler('martha', 'marhta') = %.3f\n", sim)

	// Cosine similarity
	cos, _ := edlib.Quick().Cosine("hello world", "hello there")
	fmt.Printf("Cosine('hello world', 'hello there') = %.3f\n\n", cos)
}

func fluentExamples() {
	fmt.Println("2. Fluent API Examples:")
	fmt.Println("-----------------------")

	// Basic comparison
	result, _ := edlib.New().
		Using(algorithms.Levenshtein).
		Compare("hello", "helo")
	fmt.Printf("Distance: %d, Similarity: %.3f\n", result.Distance, result.Similarity)

	// With threshold
	result2, _ := edlib.New().
		Using(algorithms.JaroWinkler).
		WithThreshold(0.8).
		WithCaseSensitive(false).
		Compare("Hello", "hello")
	fmt.Printf("Case-insensitive similarity: %.3f\n", result2.Similarity)

	// With custom weights
	result3, _ := edlib.New().
		Using(algorithms.Levenshtein).
		WithWeights(edlib.CostWeights{
			Insertion:    1,
			Deletion:     1,
			Substitution: 2,
		}).
		Compare("hello", "hallo")
	fmt.Printf("With custom weights, distance: %d\n\n", result3.Distance)
}

func searchExamples() {
	fmt.Println("3. Fuzzy Search Examples:")
	fmt.Println("-------------------------")

	targets := []string{
		"apple", "application", "apply", "appreciate",
		"banana", "bandana", "cabana",
		"cherry", "merry", "berry",
	}

	// Find best match
	best, _ := edlib.New().
		Search().
		WithAlgorithm(algorithms.Levenshtein).
		FindBest("aplication", targets)
	fmt.Printf("Best match for 'aplication': %s (score: %.3f)\n", best.Value, best.Score)

	// Find top 3 matches
	top3, _ := edlib.New().
		Search().
		WithAlgorithm(algorithms.JaroWinkler).
		FindTopN("appl", targets, 3)
	fmt.Println("\nTop 3 matches for 'appl':")
	for _, result := range top3 {
		fmt.Printf("  %d. %s (score: %.3f)\n", result.Metadata.Rank, result.Value, result.Score)
	}

	// Find with threshold
	matches, _ := edlib.New().
		Search().
		WithAlgorithm(algorithms.Levenshtein).
		WithThreshold(0.7).
		FindAll("banana", targets)
	fmt.Printf("\nMatches for 'banana' with threshold 0.7: %d results\n\n", len(matches))
}

func batchExamples() {
	fmt.Println("4. Batch Processing Examples:")
	fmt.Println("------------------------------")

	pairs := []edlib.StringPair{
		{First: "hello", Second: "helo"},
		{First: "world", Second: "wrold"},
		{First: "test", Second: "testing"},
	}

	results, _ := edlib.New().
		Using(algorithms.Levenshtein).
		Batch(pairs)

	fmt.Println("Batch results:")
	for i, result := range results {
		fmt.Printf("  %d. '%s' vs '%s': distance=%d, similarity=%.3f\n",
			i+1, pairs[i].First, pairs[i].Second, result.Distance, result.Similarity)
	}
	fmt.Println()
}

func advancedExamples() {
	fmt.Println("5. Advanced Features:")
	fmt.Println("---------------------")

	// With metrics
	result, _ := edlib.New().
		Using(algorithms.Levenshtein).
		WithMetrics(true).
		WithMemoryOptimization(true).
		WithEarlyTermination(true).
		Compare("hello world", "hello there")

	fmt.Printf("Distance: %d\n", result.Distance)
	if result.Metrics != nil {
		fmt.Printf("Execution time: %v\n", result.Metrics.ExecutionTime)
		fmt.Printf("Memory allocated: %d bytes\n", result.Metrics.MemoryAllocated)
		fmt.Printf("Optimizations used: %v\n", result.Metrics.OptimizationsUsed)
	}

	// Unicode support
	unicodeDist, _ := edlib.Quick().Levenshtein("こんにちは", "こんにちわ")
	fmt.Printf("\nUnicode distance: %d\n", unicodeDist)
}

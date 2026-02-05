package edlib_test

import (
	"context"
	"testing"
	"time"

	"github.com/hbollon/go-edlib/v2"
	"github.com/hbollon/go-edlib/v2/algorithms"
)

func TestQuickAPI(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected int
	}{
		{"identical", "hello", "hello", 0},
		{"one char diff", "hello", "hallo", 1},
		{"insertion", "hello", "helllo", 1},
		{"deletion", "hello", "helo", 1},
		{"empty first", "", "hello", 5},
		{"empty second", "hello", "", 5},
		{"both empty", "", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dist, err := edlib.Quick().Levenshtein(tt.s1, tt.s2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dist != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, dist)
			}
		})
	}
}

func TestFluentAPI(t *testing.T) {
	t.Run("basic comparison", func(t *testing.T) {
		result, err := edlib.New().
			Using(algorithms.Levenshtein).
			Compare("kitten", "sitting")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Distance != 3 {
			t.Errorf("expected distance 3, got %d", result.Distance)
		}

		if result.Algorithm != algorithms.Levenshtein {
			t.Errorf("expected Levenshtein algorithm, got %v", result.Algorithm)
		}
	})

	t.Run("with threshold", func(t *testing.T) {
		result, err := edlib.New().
			Using(algorithms.Levenshtein).
			WithThreshold(0.8).
			Compare("hello", "helo")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Similarity < 0.8 {
			t.Errorf("expected similarity >= 0.8, got %f", result.Similarity)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		result, err := edlib.New().
			Using(algorithms.Levenshtein).
			WithCaseSensitive(false).
			Compare("Hello", "hello")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Distance != 0 {
			t.Errorf("expected distance 0 for case-insensitive comparison, got %d", result.Distance)
		}
	})
}

func TestJaroWinkler(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		minScore float64
	}{
		{"identical", "hello", "hello", 1.0},
		{"similar", "martha", "marhta", 0.9},
		{"prefix bonus", "dixon", "dicksonx", 0.8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim, err := edlib.Quick().JaroWinkler(tt.s1, tt.s2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if sim < tt.minScore {
				t.Errorf("expected similarity >= %f, got %f", tt.minScore, sim)
			}
		})
	}
}

func TestCosineSimilarity(t *testing.T) {
	sim, err := edlib.Quick().Cosine("hello world", "hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sim < 0.9999 || sim > 1.0001 {
		t.Errorf("expected similarity ~1.0 for identical strings, got %f", sim)
	}

	sim2, err := edlib.Quick().Cosine("hello world", "world hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sim2 < 0.15 {
		t.Errorf("expected similarity >= 0.15 for similar strings with some common n-grams, got %f", sim2)
	}
}


func TestSearch(t *testing.T) {
	targets := []string{"test", "tester", "testing", "best", "rest"}

	t.Run("find best", func(t *testing.T) {
		result, err := edlib.New().
			Search().
			WithAlgorithm(algorithms.Levenshtein).
			FindBest("tset", targets)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Value != "test" {
			t.Errorf("expected 'test', got %s", result.Value)
		}
	})

	t.Run("find with threshold", func(t *testing.T) {
		results, err := edlib.New().
			Search().
			WithAlgorithm(algorithms.Levenshtein).
			WithThreshold(0.7).
			FindAll("test", targets)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) < 2 {
			t.Errorf("expected at least 2 results, got %d", len(results))
		}

		// Check all results meet threshold
		for _, r := range results {
			if r.Score < 0.7 {
				t.Errorf("result %s has score %f below threshold 0.7", r.Value, r.Score)
			}
		}
	})

	t.Run("find top N", func(t *testing.T) {
		results, err := edlib.New().
			Search().
			WithAlgorithm(algorithms.Levenshtein).
			FindTopN("test", targets, 3)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 3 {
			t.Errorf("expected 3 results, got %d", len(results))
		}

		// Check results are ranked
		for i, r := range results {
			if r.Metadata.Rank != i+1 {
				t.Errorf("expected rank %d, got %d", i+1, r.Metadata.Rank)
			}
		}
	})
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(2 * time.Nanosecond)

	_, err := edlib.New().
		WithContext(ctx).
		Using(algorithms.Levenshtein).
		Compare("hello", "world")

	if err == nil {
		t.Error("expected context cancellation error")
	}
}

func TestUnicodeSupport(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
	}{
		{"emoji", "hello 👋", "hello 👋"},
		{"japanese", "こんにちは", "こんにちは"},
		{"mixed", "Hello 世界", "Hello 世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dist, err := edlib.Quick().Levenshtein(tt.s1, tt.s2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dist != 0 {
				t.Errorf("expected distance 0 for identical Unicode strings, got %d", dist)
			}
		})
	}
}

func TestBatchProcessing(t *testing.T) {
	pairs := []edlib.StringPair{
		{First: "hello", Second: "helo"},
		{First: "world", Second: "wrold"},
		{First: "test", Second: "test"},
	}

	results, err := edlib.New().
		Using(algorithms.Levenshtein).
		Batch(pairs)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != len(pairs) {
		t.Errorf("expected %d results, got %d", len(pairs), len(results))
	}

	// Third pair should be identical
	if results[2].Distance != 0 {
		t.Errorf("expected distance 0 for identical pair, got %d", results[2].Distance)
	}
}

func TestWithMetrics(t *testing.T) {
	result, err := edlib.New().
		Using(algorithms.Levenshtein).
		WithMetrics(true).
		WithMemoryOptimization(true).
		Compare("hello", "world")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Metrics == nil {
		t.Error("expected metrics to be populated")
	}

	// ExecutionTime might be 0 on very fast systems, so we just check that metrics are collected
	if result.Metrics != nil && result.Metrics.ExecutionTime < 0 {
		t.Error("expected non-negative execution time")
	}
}


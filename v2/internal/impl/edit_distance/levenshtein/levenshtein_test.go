package levenshtein

import (
	"context"
	"testing"
	"time"

	"github.com/hbollon/go-edlib/v2/internal/core"
)

func TestLevenshteinDistance(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"First arg empty", args{"", "abcde"}, 5},
		{"Second arg empty", args{"abcde", ""}, 5},
		{"Same args", args{"abcde", "abcde"}, 0},
		{"ab/aa", args{"ab", "aa"}, 1},
		{"ab/ba", args{"ab", "ba"}, 2},
		{"ab/aaa", args{"ab", "aaa"}, 2},
		{"bbb/a", args{"bbb", "a"}, 3},
		{"kitten/sitting", args{"kitten", "sitting"}, 3},
		{"distance/difference", args{"distance", "difference"}, 5},
		{"a cat/an abct", args{"a cat", "an abct"}, 4},
		{"こにんち/こんにちは", args{"こにんち", "こんにちは"}, 3}, // "Hello" in Japanese
		{"🙂😄🙂😄/😄🙂😄🙂", args{"🙂😄🙂😄", "😄🙂😄🙂"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calc.Distance(tt.args.str1, tt.args.str2); got != tt.want {
				t.Errorf("Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevenshteinSimilarity(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"Same strings", args{"hello", "hello"}, 1.0},
		{"Empty strings", args{"", ""}, 1.0},
		{"Completely different", args{"abc", "xyz"}, 0.0},
		{"One char difference", args{"hello", "hallo"}, 0.8}, // 1/5 = 0.8
		{"kitten/sitting", args{"kitten", "sitting"}, 0.5714285714285714}, // 3/7
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Similarity(tt.args.str1, tt.args.str2)
			if diff := got - tt.want; diff > 0.0001 || diff < -0.0001 {
				t.Errorf("Similarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevenshteinWithCustomWeights(t *testing.T) {
	// Create calculator with custom weights (insertion=1, deletion=2, substitution=3)
	calc := &Calculator{
		weights: weights{
			insertion:    1,
			deletion:     2,
			substitution: 3,
		},
	}

	tests := []struct {
		name string
		s1   string
		s2   string
		want int
	}{
		{"Empty to abc (3 insertions)", "", "abc", 3},    // 3 insertions * 1
		{"abc to empty (3 deletions)", "abc", "", 6},    // 3 deletions * 2
		{"a to b (1 substitution)", "a", "b", 3},        // 1 substitution * 3
		{"ab to ba", "ab", "ba", 3},                     // 1 substitution * 3 (or 2 insertions/deletions)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calc.Distance(tt.s1, tt.s2); got != tt.want {
				t.Errorf("Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevenshteinCaseInsensitive(t *testing.T) {
	calc := &Calculator{
		caseInsensitive: true,
		weights: weights{
			insertion:    1,
			deletion:     1,
			substitution: 1,
		},
	}

	tests := []struct {
		name string
		s1   string
		s2   string
		want int
	}{
		{"HELLO/hello", "HELLO", "hello", 0},
		{"Hello/hello", "Hello", "hello", 0},
		{"HELLO/WORLD", "HELLO", "WORLD", 4},
		{"Hello/World", "Hello", "World", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calc.Distance(tt.s1, tt.s2); got != tt.want {
				t.Errorf("Distance() with case insensitive = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevenshteinWithMemoryPool(t *testing.T) {
	pool := core.NewMemoryPool(10000)
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}
	calc.WithMemoryPool(pool)

	// Run some calculations to test pool usage
	tests := []struct {
		s1   string
		s2   string
		want int
	}{
		{"hello", "world", 4},
		{"kitten", "sitting", 3},
		{"distance", "difference", 5},
	}

	for _, tt := range tests {
		if got := calc.Distance(tt.s1, tt.s2); got != tt.want {
			t.Errorf("Distance() with pool = %v, want %v", got, tt.want)
		}
	}

	// Verify pool statistics exist
	hits, misses := pool.Stats()
	if hits == 0 && misses == 0 {
		t.Error("Expected pool to have recorded statistics")
	}
}

func TestLevenshteinWithContext(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	t.Run("Normal execution", func(t *testing.T) {
		ctx := context.Background()
		dist, err := calc.DistanceWithContext(ctx, "hello", "world")
		if err != nil {
			t.Errorf("DistanceWithContext() error = %v", err)
		}
		if dist != 4 {
			t.Errorf("DistanceWithContext() = %v, want 4", dist)
		}
	})

	t.Run("Cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		_, err := calc.DistanceWithContext(ctx, "hello", "world")
		if err == nil {
			t.Error("DistanceWithContext() expected error with cancelled context")
		}
	})

	t.Run("Timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		time.Sleep(10 * time.Millisecond) // Ensure timeout
		
		_, err := calc.DistanceWithContext(ctx, "hello", "world")
		if err == nil {
			t.Error("DistanceWithContext() expected error with timeout context")
		}
	})
}

func TestLevenshteinEarlyTermination(t *testing.T) {
	calc := &Calculator{
		weights: weights{
			insertion:    1,
			deletion:     1,
			substitution: 1,
		},
		earlyTermination: true,
		maxDistance:      3,
	}

	// Test case where distance exceeds maxDistance
	dist := calc.Distance("hello", "world")
	// Distance should be >= maxDistance when early termination kicks in
	if dist < calc.maxDistance {
		t.Errorf("Expected early termination to trigger, got distance %v", dist)
	}
}

func TestLevenshteinMaxDistance(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	tests := []struct {
		name string
		s1   string
		s2   string
		want int
	}{
		{"Empty strings", "", "", 0},
		{"One empty", "abc", "", 3},
		{"Other empty", "", "xyz", 3},
		{"Same length", "abc", "xyz", 3},
		{"Different length", "ab", "abcd", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calc.MaxDistance(tt.s1, tt.s2); got != tt.want {
				t.Errorf("MaxDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevenshteinProperties(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	if calc.Name() != "Levenshtein" {
		t.Errorf("Name() = %v, want Levenshtein", calc.Name())
	}

	props := calc.Properties()
	if props.Name != "Levenshtein" {
		t.Errorf("Properties().Name = %v, want Levenshtein", props.Name)
	}
}

func BenchmarkLevenshteinDistance(b *testing.B) {
	calc, _ := New()
	s1, s2 := "sitting", "kitten"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Distance(s1, s2)
	}
}

func BenchmarkLevenshteinDistanceWithPool(b *testing.B) {
	pool := core.NewMemoryPool(10000)
	calc, _ := New()
	calc.WithMemoryPool(pool)
	s1, s2 := "sitting", "kitten"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Distance(s1, s2)
	}
}

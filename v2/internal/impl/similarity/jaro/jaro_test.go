package jaro

import (
	"context"
	"testing"
	"time"

	"github.com/hbollon/go-edlib/v2/internal/core"
)

func TestJaroSimilarity(t *testing.T) {
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
		{"First arg empty", args{"", "abcde"}, 0.0},
		{"Second arg empty", args{"abcde", ""}, 0.0},
		{"Both empty", args{"", ""}, 1.0},
		{"Same args", args{"abcde", "abcde"}, 1.0},
		{"No characters match", args{"abcd", "effgghh"}, 0.0},
		{"CRATE/TRACE", args{"CRATE", "TRACE"}, 0.73333335},
		{"MARTHA/MARHTA", args{"MARTHA", "MARHTA"}, 0.9444444},
		{"DIXON/DICKSONX", args{"DIXON", "DICKSONX"}, 0.76666665},
		{"jellyfish/smellyfish", args{"jellyfish", "smellyfish"}, 0.8962963},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Similarity(tt.args.str1, tt.args.str2)
			// Use tolerance for floating point comparison
			if diff := got - tt.want; diff > 0.0001 || diff < -0.0001 {
				t.Errorf("JaroSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJaroWinklerSimilarity(t *testing.T) {
	calc, err := NewWinkler()
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
		{"First arg empty", args{"", "abcde"}, 0.0},
		{"Second arg empty", args{"abcde", ""}, 0.0},
		{"Both empty", args{"", ""}, 1.0},
		{"Same args", args{"abcde", "abcde"}, 1.0},
		{"No characters match", args{"abcd", "effgghh"}, 0.0},
		{"TRACE/TRACE", args{"TRACE", "TRACE"}, 1.0},
		{"CRATE/TRACE", args{"CRATE", "TRACE"}, 0.73333335},
		{"TRATE/TRACE", args{"TRATE", "TRACE"}, 0.90666664},
		{"DIXON/DICKSONX", args{"DIXON", "DICKSONX"}, 0.81333333},
		{"jellyfish/smellyfish", args{"jellyfish", "smellyfish"}, 0.8962963}, // No common prefix
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Similarity(tt.args.str1, tt.args.str2)
			// Use tolerance for floating point comparison
			if diff := got - tt.want; diff > 0.0001 || diff < -0.0001 {
				t.Errorf("JaroWinklerSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJaroWinklerPrefixBonus(t *testing.T) {
	calc, err := NewWinkler()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	// Test that strings with common prefix get higher similarity
	// TRATE vs TRACE have common prefix "TR" (2 chars)
	withPrefix := calc.Similarity("TRATE", "TRACE")
	
	// CRATE vs TRACE have no common prefix
	withoutPrefix := calc.Similarity("CRATE", "TRACE")

	if withPrefix <= withoutPrefix {
		t.Errorf("Expected prefix bonus: withPrefix=%v should be > withoutPrefix=%v", 
			withPrefix, withoutPrefix)
	}
}

func TestJaroCaseInsensitive(t *testing.T) {
	calc := &Calculator{
		useWinkler:      false,
		caseInsensitive: true,
	}

	tests := []struct {
		name string
		s1   string
		s2   string
		want float64
	}{
		{"MARTHA/martha", "MARTHA", "martha", 1.0},
		{"Martha/martha", "Martha", "martha", 1.0},
		{"DIXON/dixon", "DIXON", "dixon", 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Similarity(tt.s1, tt.s2)
			if diff := got - tt.want; diff > 0.0001 || diff < -0.0001 {
				t.Errorf("Similarity() with case insensitive = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJaroWithMemoryPool(t *testing.T) {
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
		min  float64
	}{
		{"MARTHA", "MARHTA", 0.9},
		{"DIXON", "DICKSONX", 0.7},
		{"jellyfish", "smellyfish", 0.8},
	}

	for _, tt := range tests {
		if got := calc.Similarity(tt.s1, tt.s2); got < tt.min {
			t.Errorf("Similarity() with pool = %v, want >= %v", got, tt.min)
		}
	}

	// Verify pool statistics exist
	hits, misses := pool.Stats()
	if hits == 0 && misses == 0 {
		t.Error("Expected pool to have recorded statistics")
	}
}

func TestJaroDistance(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	// Test that Distance is the inverse of Similarity
	s1, s2 := "hello", "world"
	similarity := calc.Similarity(s1, s2)
	distance := calc.Distance(s1, s2)
	
	expectedDist := int((1.0 - similarity) * 100)
	if distance != expectedDist {
		t.Errorf("Distance() = %v, want %v", distance, expectedDist)
	}
}

func TestJaroWithContext(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	t.Run("Normal execution", func(t *testing.T) {
		ctx := context.Background()
		sim, err := calc.SimilarityWithContext(ctx, "hello", "world")
		if err != nil {
			t.Errorf("SimilarityWithContext() error = %v", err)
		}
		if sim < 0 || sim > 1 {
			t.Errorf("SimilarityWithContext() = %v, want value in [0, 1]", sim)
		}
	})

	t.Run("Cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		_, err := calc.SimilarityWithContext(ctx, "hello", "world")
		if err == nil {
			t.Error("SimilarityWithContext() expected error with cancelled context")
		}
	})

	t.Run("Timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		time.Sleep(10 * time.Millisecond) // Ensure timeout
		
		_, err := calc.SimilarityWithContext(ctx, "hello", "world")
		if err == nil {
			t.Error("SimilarityWithContext() expected error with timeout context")
		}
	})
}

func TestJaroMaxDistance(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	// MaxDistance for similarity algorithms should return 100 (representing 100% dissimilarity)
	maxDist := calc.MaxDistance("hello", "world")
	if maxDist != 100 {
		t.Errorf("MaxDistance() = %v, want 100", maxDist)
	}
}

func TestJaroProperties(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	if calc.Name() != "Jaro" {
		t.Errorf("Name() = %v, want Jaro", calc.Name())
	}

	props := calc.Properties()
	if props.Name != "Jaro" {
		t.Errorf("Properties().Name = %v, want Jaro", props.Name)
	}
}

func TestJaroWinklerProperties(t *testing.T) {
	calc, err := NewWinkler()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	if calc.Name() != "JaroWinkler" {
		t.Errorf("Name() = %v, want JaroWinkler", calc.Name())
	}

	props := calc.Properties()
	if props.Name != "JaroWinkler" {
		t.Errorf("Properties().Name = %v, want JaroWinkler", props.Name)
	}
}

func BenchmarkJaroSimilarity(b *testing.B) {
	calc, _ := New()
	s1, s2 := "MARTHA", "MARHTA"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Similarity(s1, s2)
	}
}

func BenchmarkJaroWinklerSimilarity(b *testing.B) {
	calc, _ := NewWinkler()
	s1, s2 := "DIXON", "DICKSONX"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Similarity(s1, s2)
	}
}

func BenchmarkJaroWithPool(b *testing.B) {
	pool := core.NewMemoryPool(10000)
	calc, _ := New()
	calc.WithMemoryPool(pool)
	s1, s2 := "MARTHA", "MARHTA"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Similarity(s1, s2)
	}
}

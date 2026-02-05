package cosine

import (
	"context"
	"testing"
	"time"
)

func TestCosineSimilarity(t *testing.T) {
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
		{"Empty strings", args{"", ""}, 0.0},
		{"One empty", args{"hello", ""}, 0.0},
		{"Other empty", args{"", "world"}, 0.0},
		{"Completely different", args{"abc", "xyz"}, 0.0},
		{"Similar strings", args{"hello", "hallo"}, 0.5}, // Approximate
		{"Anagrams", args{"listen", "silent"}, 0.2},      // Different n-gram positions
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Similarity(tt.args.str1, tt.args.str2)
			if got < 0.0 || got > 1.0 {
				t.Errorf("Similarity() = %v, want value in [0, 1]", got)
			}
			// For specific expected values, check with tolerance
			if tt.want == 1.0 && (got < 0.99 || got > 1.01) {
				t.Errorf("Similarity() = %v, want ~1.0", got)
			}
			if tt.want == 0.0 && got != 0.0 {
				t.Errorf("Similarity() = %v, want 0.0", got)
			}
		})
	}
}

func TestCosineSimilarityWithSplitOnSpace(t *testing.T) {
	calc := &Calculator{
		ngramSize:    2,
		splitOnSpace: true,
	}

	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name string
		args args
		min  float64 // Minimum expected similarity
	}{
		{"Same sentence", args{"I love horror movies", "I love horror movies"}, 1.0},
		{"Different sentences", args{"Radiohead", "Carly Rae Jepsen"}, 0.0},
		{"Partial match", args{"I love horror movies", "Lights out is a horror movie"}, 0.15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Similarity(tt.args.str1, tt.args.str2)
			if tt.min == 1.0 {
				if got < 0.99 {
					t.Errorf("Similarity() = %v, want ~1.0", got)
				}
			} else if tt.min == 0.0 {
				if got != 0.0 {
					t.Errorf("Similarity() = %v, want 0.0", got)
				}
			} else if got < tt.min {
				t.Errorf("Similarity() = %v, want >= %v", got, tt.min)
			}
		})
	}
}

func TestCosineSimilarityWithNgrams(t *testing.T) {
	// Test with different n-gram sizes
	tests := []struct {
		name      string
		ngramSize int
		s1        string
		s2        string
		wantSim   bool // true if similarity should be > 0
	}{
		{"2-gram same", 2, "hello", "hello", true},
		{"3-gram same", 3, "hello", "hello", true},
		{"2-gram different", 2, "abc", "xyz", false},
		{"2-gram similar", 2, "test", "best", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := &Calculator{
				ngramSize: tt.ngramSize,
			}
			got := calc.Similarity(tt.s1, tt.s2)
			
			if tt.wantSim && got == 0.0 {
				t.Errorf("Similarity() = 0.0, expected > 0 for similar strings")
			}
			if !tt.wantSim && got != 0.0 {
				t.Errorf("Similarity() = %v, expected 0.0 for different strings", got)
			}
		})
	}
}

func TestCosineCaseInsensitive(t *testing.T) {
	calc := &Calculator{
		ngramSize:       2,
		caseInsensitive: true,
	}

	tests := []struct {
		name string
		s1   string
		s2   string
		want float64
	}{
		{"HELLO/hello", "HELLO", "hello", 1.0},
		{"Hello/hello", "Hello", "hello", 1.0},
		{"HeLLo/hello", "HeLLo", "hello", 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Similarity(tt.s1, tt.s2)
			if diff := got - tt.want; diff > 0.01 || diff < -0.01 {
				t.Errorf("Similarity() with case insensitive = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCosineDistance(t *testing.T) {
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

	// Test edge cases
	tests := []struct {
		name string
		s1   string
		s2   string
		want int
	}{
		{"Same strings", "hello", "hello", 0},
		{"Empty strings", "", "", 100}, // 0 similarity = 100 distance
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Distance(tt.s1, tt.s2)
			if got != tt.want {
				t.Errorf("Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCosineWithContext(t *testing.T) {
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

func TestCosineMaxDistance(t *testing.T) {
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

func TestCosineProperties(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	if calc.Name() != "Cosine" {
		t.Errorf("Name() = %v, want Cosine", calc.Name())
	}

	props := calc.Properties()
	if props.Name != "Cosine" {
		t.Errorf("Properties().Name = %v, want Cosine", props.Name)
	}
}

func TestCosineUnicodeSupport(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	tests := []struct {
		name string
		s1   string
		s2   string
		want float64 // Exact or minimum expected
	}{
		{"Japanese identical", "こんにちは", "こんにちは", 1.0},
		{"Emojis identical", "🙂😄🙂😄", "🙂😄🙂😄", 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Similarity(tt.s1, tt.s2)
			if tt.want == 1.0 && got < 0.99 {
				t.Errorf("Similarity() = %v, want ~1.0", got)
			}
		})
	}
}

func BenchmarkCosineSimilarity(b *testing.B) {
	calc, _ := New()
	s1, s2 := "hello world", "hello there"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Similarity(s1, s2)
	}
}

func BenchmarkCosineSimilarityWithSplit(b *testing.B) {
	calc := &Calculator{
		ngramSize:    2,
		splitOnSpace: true,
	}
	s1, s2 := "I love horror movies", "Lights out is a horror movie"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Similarity(s1, s2)
	}
}

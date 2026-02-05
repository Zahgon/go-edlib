package damerau

import (
	"context"
	"testing"
	"time"

	"github.com/hbollon/go-edlib/v2/algorithms"
)

func TestOSADamerauLevenshteinDistance(t *testing.T) {
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
		{"ab/ba", args{"ab", "ba"}, 1},
		{"ab/aaa", args{"ab", "aaa"}, 2},
		{"bbb/a", args{"bbb", "a"}, 3},
		{"ca/abc", args{"ca", "abc"}, 3},
		{"a cat/an abct", args{"a cat", "an abct"}, 4},
		{"dixon/dicksonx", args{"dixon", "dicksonx"}, 4},
		{"jellyfish/smellyfish", args{"jellyfish", "smellyfish"}, 2},
		{"こにんち/こんにちは", args{"こにんち", "こんにちは"}, 2}, // "Hello" in Japanese
		{"🙂😄🙂😄/😄🙂😄🙂", args{"🙂😄🙂😄", "😄🙂😄🙂"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := NewOSA()
			if err != nil {
				t.Fatalf("NewOSA() error = %v", err)
			}
			if got := calc.Distance(tt.args.str1, tt.args.str2); got != tt.want {
				t.Errorf("OSADamerauLevenshtein Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDamerauLevenshteinDistance(t *testing.T) {
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
		{"ab/ba", args{"ab", "ba"}, 1},
		{"ab/aaa", args{"ab", "aaa"}, 2},
		{"bbb/a", args{"bbb", "a"}, 3},
		{"ca/abc", args{"ca", "abc"}, 2},
		{"a cat/an abct", args{"a cat", "an abct"}, 3},
		{"dixon/dicksonx", args{"dixon", "dicksonx"}, 4},
		{"jellyfish/smellyfish", args{"jellyfish", "smellyfish"}, 2},
		{"こにんち/こんにちは", args{"こにんち", "こんにちは"}, 2}, // "Hello" in Japanese
		{"🙂😄🙂😄/😄🙂😄🙂", args{"🙂😄🙂😄", "😄🙂😄🙂"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			if got := calc.Distance(tt.args.str1, tt.args.str2); got != tt.want {
				t.Errorf("DamerauLevenshtein Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDamerauLevenshteinCaseInsensitive(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	
	// Case sensitive (default)
	dist1 := calc.Distance("Hello", "hello")
	if dist1 == 0 {
		t.Errorf("Case-sensitive distance should not be 0 for 'Hello' vs 'hello'")
	}

	// Case insensitive
	if err := calc.SetCaseInsensitive(true); err != nil {
		t.Fatalf("SetCaseInsensitive() error = %v", err)
	}
	dist2 := calc.Distance("Hello", "hello")
	if dist2 != 0 {
		t.Errorf("Case-insensitive distance = %v, want 0", dist2)
	}
}

func TestDamerauLevenshteinSimilarity(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	tests := []struct {
		name string
		s1   string
		s2   string
		want float64
	}{
		{"Identical strings", "hello", "hello", 1.0},
		{"Empty strings", "", "", 1.0},
		{"Completely different", "abc", "xyz", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := calc.Similarity(tt.s1, tt.s2)
			if sim != tt.want {
				t.Errorf("Similarity() = %v, want %v", sim, tt.want)
			}
		})
	}
}

func TestDamerauLevenshteinWithContext(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("Normal execution", func(t *testing.T) {
		ctx := context.Background()
		dist, err := calc.DistanceWithContext(ctx, "hello", "hallo")
		if err != nil {
			t.Errorf("DistanceWithContext() error = %v", err)
		}
		if dist != 1 {
			t.Errorf("DistanceWithContext() = %v, want 1", dist)
		}
	})

	t.Run("Cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := calc.DistanceWithContext(ctx, "hello", "world")
		if err == nil {
			t.Error("DistanceWithContext() should return error for cancelled context")
		}
	})

	t.Run("Timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		time.Sleep(2 * time.Nanosecond)
		_, err := calc.DistanceWithContext(ctx, "hello", "world")
		if err == nil {
			t.Error("DistanceWithContext() should return error for timed out context")
		}
	})
}

func TestDamerauLevenshteinProperties(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if calc.Name() != "DamerauLevenshtein" {
		t.Errorf("Name() = %v, want DamerauLevenshtein", calc.Name())
	}

	props := calc.Properties()
	if props.Type != algorithms.EditDistanceType {
		t.Errorf("Properties().Type = %v, want %v", props.Type, algorithms.EditDistanceType)
	}
	if !props.SupportsUnicode {
		t.Error("Properties().SupportsUnicode should be true")
	}
}

func TestOSAProperties(t *testing.T) {
	calc, err := NewOSA()
	if err != nil {
		t.Fatalf("NewOSA() error = %v", err)
	}

	if calc.Name() != "OSADamerauLevenshtein" {
		t.Errorf("Name() = %v, want OSADamerauLevenshtein", calc.Name())
	}

	props := calc.Properties()
	if props.Type != algorithms.EditDistanceType {
		t.Errorf("Properties().Type = %v, want %v", props.Type, algorithms.EditDistanceType)
	}
}

func TestMaxDistance(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	tests := []struct {
		name string
		s1   string
		s2   string
		want int
	}{
		{"Same length", "hello", "world", 5},
		{"Different length", "hi", "hello", 5},
		{"Empty first", "", "hello", 5},
		{"Empty second", "hello", "", 5},
		{"Both empty", "", "", 0},
		{"Unicode", "こんにちは", "こん", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			max := calc.MaxDistance(tt.s1, tt.s2)
			if max != tt.want {
				t.Errorf("MaxDistance() = %v, want %v", max, tt.want)
			}
		})
	}
}

func BenchmarkOSADamerauLevenshtein(b *testing.B) {
	calc, _ := NewOSA()
	s1 := "kitten"
	s2 := "sitting"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Distance(s1, s2)
	}
}

func BenchmarkDamerauLevenshtein(b *testing.B) {
	calc, _ := New()
	s1 := "kitten"
	s2 := "sitting"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Distance(s1, s2)
	}
}

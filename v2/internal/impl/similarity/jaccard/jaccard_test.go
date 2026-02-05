package jaccard

import (
	"context"
	"testing"
	"time"
)

func TestJaccardSimilarity(t *testing.T) {
	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{"Jaccard sim 1", args{"Radiohead", "Carly Rae Jepsen"}, 0.0},
		{"Jaccard sim 2", args{"I love horror movies", "Lights out is a horror movie"}, 1.0 / 9.0},
		{"Jaccard sim 3", args{"love horror movies", "Lights out horror movie"}, 1.0 / 6.0},
		{"Jaccard sim 4", args{"私の名前はジョンです", "私の名前はジョン・ドゥです"}, 0.0},
		{"Jaccard sim 5", args{"🙂😄🙂😄 😄🙂😄", "🙂😄🙂😄 😄🙂😄 🙂😄🙂"}, 2.0 / 3.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			calc.SetNgramSize(0) // Split on whitespace
			if got := float32(calc.Similarity(tt.args.str1, tt.args.str2)); got != tt.want {
				t.Errorf("Jaccard Similarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJaccardShingleSimilarity(t *testing.T) {
	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{"Jaccard shingle sim 1", args{"Radiohead", "Carly Rae Jepsen"}, 0.04761905},
		{"Jaccard shingle sim 2", args{"I love horror movies", "Lights out is a horror movie"}, 0.3548387},
		{"Jaccard shingle sim 3", args{"love horror movies", "Lights out horror movie"}, 0.44},
		{"Jaccard shingle sim 4", args{"私の名前はジョンです", "私の名前はジョン・ドゥです"}, 0.61538464},
		{"Jaccard shingle sim 5", args{"🙂😄🙂😄 😄🙂😄", "🙂😄🙂😄 😄🙂😄 🙂😄🙂"}, 0.8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			calc.SetNgramSize(2) // Use 2-gram shingles
			got := float32(calc.Similarity(tt.args.str1, tt.args.str2))
			if diff := got - tt.want; diff > 0.0001 || diff < -0.0001 {
				t.Errorf("Jaccard Similarity() with shingle 2 = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJaccardWithContext(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("Normal execution", func(t *testing.T) {
		ctx := context.Background()
		sim, err := calc.SimilarityWithContext(ctx, "hello", "hello")
		if err != nil {
			t.Errorf("SimilarityWithContext() error = %v", err)
		}
		if sim != 1.0 {
			t.Errorf("SimilarityWithContext() = %v, want 1.0", sim)
		}
	})

	t.Run("Cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := calc.SimilarityWithContext(ctx, "hello", "world")
		if err == nil {
			t.Error("SimilarityWithContext() should return error for cancelled context")
		}
	})

	t.Run("Timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		time.Sleep(2 * time.Nanosecond)
		_, err := calc.SimilarityWithContext(ctx, "hello", "world")
		if err == nil {
			t.Error("SimilarityWithContext() should return error for timed out context")
		}
	})
}

func BenchmarkJaccard(b *testing.B) {
	calc, _ := New()
	calc.SetNgramSize(2)
	s1 := "I love horror movies"
	s2 := "Lights out is a horror movie"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Similarity(s1, s2)
	}
}

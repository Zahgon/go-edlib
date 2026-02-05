package edlib_test

import (
	"testing"

	"github.com/hbollon/go-edlib/v2"
	"github.com/hbollon/go-edlib/v2/algorithms"
)

var (
	// Test strings of varying lengths
	shortStr1  = "hello"
	shortStr2  = "helo"
	mediumStr1 = "The quick brown fox jumps over the lazy dog"
	mediumStr2 = "The qwick browm fox jumped over a lazy dogs"
	longStr1   = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua"
	longStr2   = "Lorem ipsom dolor sit amot, consectetur adipiscing elit, sed do euismod tempor incididunt ut labore et dolure magna aliqua"

	result int
	simResult float64
)

func BenchmarkLevenshtein(b *testing.B) {
	benchmarks := []struct {
		name string
		s1   string
		s2   string
	}{
		{"short", shortStr1, shortStr2},
		{"medium", mediumStr1, mediumStr2},
		{"long", longStr1, longStr2},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var r int
			builder := edlib.New().Using(algorithms.Levenshtein)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r, _ = builder.Distance(bm.s1, bm.s2)
			}
			result = r
		})

		b.Run(bm.name+"_with_pool", func(b *testing.B) {
			var r int
			builder := edlib.New().
				Using(algorithms.Levenshtein).
				WithMemoryOptimization(true)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r, _ = builder.Distance(bm.s1, bm.s2)
			}
			result = r
		})

		b.Run(bm.name+"_with_early_term", func(b *testing.B) {
			var r int
			builder := edlib.New().
				Using(algorithms.Levenshtein).
				WithEarlyTermination(true).
				WithThreshold(0.5)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r, _ = builder.Distance(bm.s1, bm.s2)
			}
			result = r
		})
	}
}

func BenchmarkJaroWinkler(b *testing.B) {
	benchmarks := []struct {
		name string
		s1   string
		s2   string
	}{
		{"short", shortStr1, shortStr2},
		{"medium", mediumStr1, mediumStr2},
		{"long", longStr1, longStr2},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var r float64
			builder := edlib.New().Using(algorithms.JaroWinkler)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r, _ = builder.Similarity(bm.s1, bm.s2)
			}
			simResult = r
		})
	}
}

func BenchmarkCosine(b *testing.B) {
	benchmarks := []struct {
		name string
		s1   string
		s2   string
	}{
		{"short", shortStr1, shortStr2},
		{"medium", mediumStr1, mediumStr2},
		{"long", longStr1, longStr2},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var r float64
			builder := edlib.New().Using(algorithms.Cosine)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r, _ = builder.Similarity(bm.s1, bm.s2)
			}
			simResult = r
		})
	}
}

func BenchmarkSearch(b *testing.B) {
	targets := []string{
		"apple", "application", "apply", "approach", "appreciate",
		"banana", "band", "bank", "banner", "bar",
		"cat", "catch", "category", "cathedral", "cattle",
		"dog", "doctor", "document", "dollar", "domain",
		"elephant", "element", "elevator", "eleven", "email",
	}

	b.Run("find_best", func(b *testing.B) {
		builder := edlib.New().Search().WithAlgorithm(algorithms.Levenshtein)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = builder.FindBest("aplication", targets)
		}
	})

	b.Run("find_top_5", func(b *testing.B) {
		builder := edlib.New().Search().WithAlgorithm(algorithms.Levenshtein)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = builder.FindTopN("aplication", targets, 5)
		}
	})

	b.Run("find_with_threshold", func(b *testing.B) {
		builder := edlib.New().Search().
			WithAlgorithm(algorithms.Levenshtein).
			WithThreshold(0.7)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = builder.FindAll("aplication", targets)
		}
	})
}

func BenchmarkBatchProcessing(b *testing.B) {
	pairs := make([]edlib.StringPair, 100)
	for i := 0; i < 100; i++ {
		pairs[i] = edlib.StringPair{
			First:  mediumStr1,
			Second: mediumStr2,
		}
	}

	b.Run("sequential", func(b *testing.B) {
		builder := edlib.New().Using(algorithms.Levenshtein)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = builder.Batch(pairs)
		}
	})

	b.Run("parallel", func(b *testing.B) {
		builder := edlib.New().
			Using(algorithms.Levenshtein).
			WithParallel(4)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = builder.Batch(pairs)
		}
	})
}

func BenchmarkQuickAPI(b *testing.B) {
	quick := edlib.Quick()

	b.Run("levenshtein", func(b *testing.B) {
		var r int
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r, _ = quick.Levenshtein(shortStr1, shortStr2)
		}
		result = r
	})

	b.Run("jaro_winkler", func(b *testing.B) {
		var r float64
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r, _ = quick.JaroWinkler(shortStr1, shortStr2)
		}
		simResult = r
	})

	b.Run("cosine", func(b *testing.B) {
		var r float64
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r, _ = quick.Cosine(shortStr1, shortStr2)
		}
		simResult = r
	})
}

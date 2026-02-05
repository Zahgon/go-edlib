package hamming

import (
	"context"
	"testing"
	"time"
)

func TestHammingDistance(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name       string
		args       args
		want       int
		wantErr    bool
		useContext bool
	}{
		{"Same strings", args{"aa", "aa"}, 0, false, false},
		{"One difference", args{"ab", "aa"}, 1, false, false},
		{"Two differences", args{"ab", "ba"}, 2, false, false},
		{"All different", args{"abc", "xyz"}, 3, false, false},
		{"Unequal length - error", args{"ab", "aaa"}, 0, true, true},
		{"Unequal length - error 2", args{"bbb", "a"}, 0, true, true},
		{"Unicode - emojis", args{"🙂😄🙂😄", "😄🙂😄🙂"}, 4, false, false},
		{"Unicode - Japanese", args{"こんにちは", "こんばんは"}, 2, false, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useContext {
				// Use DistanceWithContext for error checking
				got, err := calc.DistanceWithContext(context.Background(), tt.args.str1, tt.args.str2)
				if (err != nil) != tt.wantErr {
					t.Errorf("DistanceWithContext() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && got != tt.want {
					t.Errorf("DistanceWithContext() = %v, want %v", got, tt.want)
				}
			} else {
				// Use Distance (simple version)
				got := calc.Distance(tt.args.str1, tt.args.str2)
				if got != tt.want {
					t.Errorf("Distance() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestHammingSimilarity(t *testing.T) {
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
		{"One char difference", args{"hello", "hallo"}, 0.8}, // 1 diff out of 5
		{"Half different", args{"abc", "axc"}, 0.6666666666666667}, // 1 diff out of 3
		{"All different", args{"abc", "xyz"}, 0.0},
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

func TestHammingCaseInsensitive(t *testing.T) {
	calc := &Calculator{
		caseInsensitive: true,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calc.Distance(tt.s1, tt.s2); got != tt.want {
				t.Errorf("Distance() with case insensitive = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHammingWithContext(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
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
		cancel() // Cancel immediately
		
		_, err := calc.DistanceWithContext(ctx, "hello", "hallo")
		if err == nil {
			t.Error("DistanceWithContext() expected error with cancelled context")
		}
	})

	t.Run("Timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		time.Sleep(10 * time.Millisecond) // Ensure timeout
		
		_, err := calc.DistanceWithContext(ctx, "hello", "hallo")
		if err == nil {
			t.Error("DistanceWithContext() expected error with timeout context")
		}
	})

	t.Run("Unequal length error", func(t *testing.T) {
		ctx := context.Background()
		_, err := calc.DistanceWithContext(ctx, "hello", "world!")
		if err == nil {
			t.Error("DistanceWithContext() expected error for unequal length strings")
		}
	})
}

func TestHammingMaxDistance(t *testing.T) {
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
		{"Same length", "abc", "xyz", 3},
		{"Same length - longer", "abcde", "fghij", 5},
		{"Different length", "ab", "abcd", 4}, // Returns max even though undefined
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calc.MaxDistance(tt.s1, tt.s2); got != tt.want {
				t.Errorf("MaxDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHammingProperties(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	if calc.Name() != "Hamming" {
		t.Errorf("Name() = %v, want Hamming", calc.Name())
	}

	props := calc.Properties()
	if props.Name != "Hamming" {
		t.Errorf("Properties().Name = %v, want Hamming", props.Name)
	}
}

func BenchmarkHammingDistance(b *testing.B) {
	calc, _ := New()
	s1, s2 := "hello world", "hallo warld"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Distance(s1, s2)
	}
}

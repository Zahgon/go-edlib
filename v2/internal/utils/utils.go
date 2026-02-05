package utils

import (
	"runtime"
	"strings"
	"unicode"
)

// Min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min3 returns the minimum of three integers
func Min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// Min4 returns the minimum of four integers
func Min4(a, b, c, d int) int {
	return Min(Min(a, b), Min(c, d))
}

// Equal compares two rune slices for equality
func Equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ToLower converts a string to lowercase (returns original runes)
func ToLower(s string) []rune {
	runes := []rune(s)
	for i, r := range runes {
		runes[i] = unicode.ToLower(r)
	}
	return runes
}

// Normalize normalizes whitespace in a string
func Normalize(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// NumCPU returns the number of CPUs available, with a reasonable max
func NumCPU() int {
	n := runtime.NumCPU()
	if n > 32 {
		return 32 // Cap at 32 to avoid excessive goroutines
	}
	return n
}

// OptimalWorkers returns the optimal number of workers for a given input size
func OptimalWorkers(inputSize int, minPerWorker int) int {
	if inputSize < minPerWorker {
		return 1
	}
	
	maxWorkers := NumCPU()
	workers := inputSize / minPerWorker
	
	if workers > maxWorkers {
		return maxWorkers
	}
	
	if workers < 1 {
		return 1
	}
	
	return workers
}

// ClampFloat64 clamps a float64 value between min and max
func ClampFloat64(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// RunesToString converts a rune slice to a string efficiently
func RunesToString(runes []rune) string {
	if len(runes) == 0 {
		return ""
	}
	return string(runes)
}

// StringToRunes converts a string to a rune slice
func StringToRunes(s string) []rune {
	if s == "" {
		return nil
	}
	return []rune(s)
}

// Shingle finds the k-gram of a string for a given k
// Takes a string and an integer as parameters and returns a map with counts.
// Returns an empty map if the string is empty or if k is 0
func Shingle(s string, k int) map[string]int {
	m := make(map[string]int)
	if s != "" && k != 0 {
		runeS := []rune(s)

		for i := 0; i < len(runeS)-k+1; i++ {
			m[string(runeS[i:i+k])]++
		}
	}
	return m
}

// ShingleSlice finds the k-gram of a string for a given k
// Takes a string and an integer as parameters and returns a slice of unique k-grams.
// Returns an empty slice if the string is empty or if k is 0
func ShingleSlice(s string, k int) []string {
	var out []string
	m := make(map[string]int)
	if s != "" && k != 0 {
		runeS := []rune(s)
		for i := 0; i < len(runeS)-k+1; i++ {
			m[string(runeS[i:i+k])]++
		}
		for k := range m {
			out = append(out, k)
		}
	}
	return out
}

// StringHashMap is HashMap substitute for string (used as a set)
type StringHashMap map[string]struct{}

// AddAll adds all elements from one StringHashMap to another
func (m StringHashMap) AddAll(srcMap StringHashMap) {
	for key := range srcMap {
		m[key] = struct{}{}
	}
}

// ToArray converts a StringHashMap to a string array
func (m StringHashMap) ToArray() []string {
	arr := make([]string, 0, len(m))
	for key := range m {
		arr = append(arr, key)
	}
	return arr
}

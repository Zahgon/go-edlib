package edlib

import (
	"github.com/hbollon/go-edlib/internal/utils"
)

// LCS takes two strings and compute their LCS(Longuest Common Subsequence)
func LCS(str1, str2 string) int {
	_ = "STUB: not implemented"
	// Convert strings to rune array to handle no-ASCII characters
	return 0
}

// Return computed lcs matrix
func lcsProcess(runeStr1, runeStr2 []rune) [][]int {
	_ = "STUB: not implemented"
	// 2D Array that will contain str1 and str2 LCS
	return nil
}

// LCSBacktrack returns all choices taken during LCS process
func LCSBacktrack(str1, str2 string) (string, error) { _ = "STUB: not implemented"; return "", nil }

func processLCSBacktrack(str1, str2 string, lcsMatrix [][]int, m, n int) string {
	_ = "STUB: not implemented"
	// Convert strings to rune array to handle no-ASCII characters
	return ""
}

// LCSBacktrackAll returns an array containing all common substrings between str1 and str2
func LCSBacktrackAll(str1, str2 string) ([]string, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

func processLCSBacktrackAll(str1, str2 string, lcsMatrix [][]int, m, n int) utils.StringHashMap {
	_ = "STUB: not implemented"
	// Convert strings to rune array to handle no-ASCII characters
	return *new(utils.StringHashMap)
}

// Map containing all commons substrings (Hash set builded from map)

// LCSDiff will backtrack through the lcs matrix and return the diff between the two sequences
func LCSDiff(str1, str2 string) ([]string, error) { _ = "STUB: not implemented"; return nil, nil }

func processLCSDiff(str1 string, str2 string, lcsMatrix [][]int, m, n int) []string {
	_ = "STUB: not implemented"
	// Convert strings to rune array to handle no-ASCII characters
	return nil
}

// LCSEditDistance determines the edit distance between two strings using LCS function
// (allow only insert and delete operations)
func LCSEditDistance(str1, str2 string) int { _ = "STUB: not implemented"; return 0 }

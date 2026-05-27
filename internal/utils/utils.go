package utils

// StringHashMap is HashMap substitute for string
type StringHashMap map[string]struct{}

// Min return the smallest integer among the two in parameters
func Min(a int, b int) int { _ = "STUB: not implemented"; return 0 }

// Max return the largest integer among the two in parameters
func Max(a int, b int) int { _ = "STUB: not implemented"; return 0 }

// Equal compare two rune arrays and return if they are equals or not
func Equal(a, b []rune) bool { _ = "STUB: not implemented"; return false }

/*
	StringHashMap methods
*/

// AddAll adds all elements from one StringHashMap to another
func (m StringHashMap) AddAll(srcMap StringHashMap) { _ = "STUB: not implemented"; return }

// ToArray convert and return an StringHashMap to string array
func (m StringHashMap) ToArray() []string { _ = "STUB: not implemented"; return nil }

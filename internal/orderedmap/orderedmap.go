package orderedmap

type pair struct {
	Key   string
	Value float32
}

// OrderedMap is a slice of pairs type with string keys and float values.
// It implement sorting methods by values.
type OrderedMap []pair

// Len return length of a given OrderedMap
func (p OrderedMap) Len() int {
	_ = "STUB: not implemented"

	// Less return if a element of an OrderedMap is smaller than another
	return 0
}

func (p OrderedMap) Less(i, j int) bool { _ = "STUB: not implemented"; return false }

// Swap two members of an OrderedMap
func (p OrderedMap) Swap(i, j int) { _ = "STUB: not implemented"; return }

// ToArray export keys of an OrderedMap into a slice
func (p OrderedMap) ToArray() []string { _ = "STUB: not implemented"; return nil }

// SortByValues sort by values an OrderedMap in decreasing order
func (p OrderedMap) SortByValues() { _ = "STUB: not implemented"; return }

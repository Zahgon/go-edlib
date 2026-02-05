package algorithms

// Algorithm represents a string comparison algorithm type
type Algorithm int

// Available algorithms
const (
	// Edit distance algorithms
	Levenshtein Algorithm = iota
	DamerauLevenshtein
	OSADamerauLevenshtein
	Hamming
	LCS

	// Similarity algorithms
	Jaro
	JaroWinkler
	Cosine
	Jaccard
	SorensenDice
	QGram
)

// String returns the string representation of the algorithm
func (a Algorithm) String() string {
	switch a {
	case Levenshtein:
		return "Levenshtein"
	case DamerauLevenshtein:
		return "DamerauLevenshtein"
	case OSADamerauLevenshtein:
		return "OSADamerauLevenshtein"
	case Hamming:
		return "Hamming"
	case LCS:
		return "LCS"
	case Jaro:
		return "Jaro"
	case JaroWinkler:
		return "JaroWinkler"
	case Cosine:
		return "Cosine"
	case Jaccard:
		return "Jaccard"
	case SorensenDice:
		return "SorensenDice"
	case QGram:
		return "QGram"
	default:
		return "Unknown"
	}
}

// AlgorithmType represents the category of an algorithm
type AlgorithmType int

const (
	EditDistanceType AlgorithmType = iota
	SequenceType
	SimilarityType
)

// Type returns the algorithm's category
func (a Algorithm) Type() AlgorithmType {
	switch a {
	case Levenshtein, DamerauLevenshtein, OSADamerauLevenshtein, Hamming:
		return EditDistanceType
	case LCS:
		return SequenceType
	case Jaro, JaroWinkler, Cosine, Jaccard, SorensenDice, QGram:
		return SimilarityType
	default:
		return EditDistanceType
	}
}

// AlgorithmProperties contains properties and capabilities of an algorithm
type AlgorithmProperties struct {
	// Name is the algorithm name
	Name string

	// Type is the algorithm category
	Type AlgorithmType

	// SupportsDistance indicates if the algorithm can compute edit distance
	SupportsDistance bool

	// SupportsSimilarity indicates if the algorithm can compute similarity score
	SupportsSimilarity bool

	// RequiresEqualLength indicates if strings must be of equal length
	RequiresEqualLength bool

	// SupportsUnicode indicates if the algorithm properly handles Unicode
	SupportsUnicode bool

	// SupportsWeights indicates if the algorithm supports custom operation weights
	SupportsWeights bool

	// SupportsTransposition indicates if the algorithm supports transposition operations
	SupportsTransposition bool

	// TimeComplexity describes the time complexity (e.g., "O(n*m)")
	TimeComplexity string

	// SpaceComplexity describes the space complexity (e.g., "O(n)")
	SpaceComplexity string

	// Description provides a brief description of the algorithm
	Description string
}

// Properties returns the properties of the algorithm
func (a Algorithm) Properties() AlgorithmProperties {
	switch a {
	case Levenshtein:
		return AlgorithmProperties{
			Name:                  "Levenshtein",
			Type:                  EditDistanceType,
			SupportsDistance:      true,
			SupportsSimilarity:    true,
			RequiresEqualLength:   false,
			SupportsUnicode:       true,
			SupportsWeights:       true,
			SupportsTransposition: false,
			TimeComplexity:        "O(n*m)",
			SpaceComplexity:       "O(min(n,m))",
			Description:           "Edit distance allowing insertions, deletions, and substitutions",
		}
	case DamerauLevenshtein:
		return AlgorithmProperties{
			Name:                  "DamerauLevenshtein",
			Type:                  EditDistanceType,
			SupportsDistance:      true,
			SupportsSimilarity:    true,
			RequiresEqualLength:   false,
			SupportsUnicode:       true,
			SupportsWeights:       false,
			SupportsTransposition: true,
			TimeComplexity:        "O(n*m)",
			SpaceComplexity:       "O(n*m)",
			Description:           "Edit distance with transpositions (full Damerau-Levenshtein)",
		}
	case OSADamerauLevenshtein:
		return AlgorithmProperties{
			Name:                  "OSADamerauLevenshtein",
			Type:                  EditDistanceType,
			SupportsDistance:      true,
			SupportsSimilarity:    true,
			RequiresEqualLength:   false,
			SupportsUnicode:       true,
			SupportsWeights:       false,
			SupportsTransposition: true,
			TimeComplexity:        "O(n*m)",
			SpaceComplexity:       "O(n*m)",
			Description:           "Optimal String Alignment distance (restricted transpositions)",
		}
	case Hamming:
		return AlgorithmProperties{
			Name:                  "Hamming",
			Type:                  EditDistanceType,
			SupportsDistance:      true,
			SupportsSimilarity:    true,
			RequiresEqualLength:   true,
			SupportsUnicode:       true,
			SupportsWeights:       false,
			SupportsTransposition: false,
			TimeComplexity:        "O(n)",
			SpaceComplexity:       "O(1)",
			Description:           "Edit distance with substitutions only (equal length strings)",
		}
	case LCS:
		return AlgorithmProperties{
			Name:                  "LCS",
			Type:                  SequenceType,
			SupportsDistance:      true,
			SupportsSimilarity:    true,
			RequiresEqualLength:   false,
			SupportsUnicode:       true,
			SupportsWeights:       false,
			SupportsTransposition: false,
			TimeComplexity:        "O(n*m)",
			SpaceComplexity:       "O(n*m)",
			Description:           "Longest Common Subsequence",
		}
	case Jaro:
		return AlgorithmProperties{
			Name:                  "Jaro",
			Type:                  SimilarityType,
			SupportsDistance:      false,
			SupportsSimilarity:    true,
			RequiresEqualLength:   false,
			SupportsUnicode:       true,
			SupportsWeights:       false,
			SupportsTransposition: false,
			TimeComplexity:        "O(n*m)",
			SpaceComplexity:       "O(n+m)",
			Description:           "Similarity metric designed for short strings",
		}
	case JaroWinkler:
		return AlgorithmProperties{
			Name:                  "JaroWinkler",
			Type:                  SimilarityType,
			SupportsDistance:      false,
			SupportsSimilarity:    true,
			RequiresEqualLength:   false,
			SupportsUnicode:       true,
			SupportsWeights:       false,
			SupportsTransposition: false,
			TimeComplexity:        "O(n*m)",
			SpaceComplexity:       "O(n+m)",
			Description:           "Jaro similarity with prefix bonus",
		}
	case Cosine:
		return AlgorithmProperties{
			Name:                  "Cosine",
			Type:                  SimilarityType,
			SupportsDistance:      false,
			SupportsSimilarity:    true,
			RequiresEqualLength:   false,
			SupportsUnicode:       true,
			SupportsWeights:       false,
			SupportsTransposition: false,
			TimeComplexity:        "O(n+m)",
			SpaceComplexity:       "O(n+m)",
			Description:           "Cosine similarity between string vectors",
		}
	case Jaccard:
		return AlgorithmProperties{
			Name:                  "Jaccard",
			Type:                  SimilarityType,
			SupportsDistance:      false,
			SupportsSimilarity:    true,
			RequiresEqualLength:   false,
			SupportsUnicode:       true,
			SupportsWeights:       false,
			SupportsTransposition: false,
			TimeComplexity:        "O(n+m)",
			SpaceComplexity:       "O(n+m)",
			Description:           "Jaccard similarity coefficient",
		}
	case SorensenDice:
		return AlgorithmProperties{
			Name:                  "SorensenDice",
			Type:                  SimilarityType,
			SupportsDistance:      false,
			SupportsSimilarity:    true,
			RequiresEqualLength:   false,
			SupportsUnicode:       true,
			SupportsWeights:       false,
			SupportsTransposition: false,
			TimeComplexity:        "O(n+m)",
			SpaceComplexity:       "O(n+m)",
			Description:           "Sørensen-Dice coefficient",
		}
	case QGram:
		return AlgorithmProperties{
			Name:                  "QGram",
			Type:                  SimilarityType,
			SupportsDistance:      true,
			SupportsSimilarity:    true,
			RequiresEqualLength:   false,
			SupportsUnicode:       true,
			SupportsWeights:       false,
			SupportsTransposition: false,
			TimeComplexity:        "O(n+m)",
			SpaceComplexity:       "O(n+m)",
			Description:           "Q-gram distance and similarity",
		}
	default:
		return AlgorithmProperties{
			Name:               a.String(),
			Type:               a.Type(),
			SupportsUnicode:    true,
			TimeComplexity:     "O(n*m)",
			SpaceComplexity:    "O(n*m)",
			Description:        "String comparison algorithm",
		}
	}
}

// IsValid checks if the algorithm is valid
func (a Algorithm) IsValid() bool {
	return a >= Levenshtein && a <= QGram
}

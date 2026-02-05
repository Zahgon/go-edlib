package core

import "errors"

// ErrNoResults is returned when no results match the search criteria
var ErrNoResults = errors.New("no results found")

package array

// PUBLIC TYPES
// ========================================================================

/*
Check if an element exist in array

Return true if found, false otherwise
*/
func IsElementInArray[T comparable](arr []T, val T) bool {
	for _, element := range arr {
		if element == val {
			return true
		}
	}

	return false
}

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

// PRIVATE FUNCTIONS
// ========================================================================

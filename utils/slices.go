package utils

// DoesSliceContain - check if a string slice contains a given string
func DoesSliceContain(dataSlice []string, toCompare string) bool {
	for _, value := range dataSlice {
		if value == toCompare {
			return true
		}
	}

	return false
}

// FIXME: move to barelyhuman/go package
// Chunk - create batches of an arbitrary typed slice into the given batch size
func Chunk[K comparable](slice []K, batchSize int) [][]K {
	var batches [][]K
	for i := 0; i < len(slice); i += batchSize {
		end := i + batchSize

		if end > len(slice) {
			end = len(slice)
		}

		batches = append(batches, slice[i:end])
	}

	return batches
}

// FIXME: move to barelyhuman/go package
// PickField - pick the value of a particular field from a slice into it's own slice
func PickField[K interface{}, V comparable](iterator []K, returner func(K) V) (result []V) {
	for _, item := range iterator {
		result = append(result, returner(item))
	}
	return
}

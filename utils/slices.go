package utils

func DoesSliceContain(dataSlice []string, toCompare string) bool {
	for _, value := range dataSlice {
		if value == toCompare {
			return true
		}
	}

	return false
}

func Chunk[K comparable](slice []K, batch int) [][]K {
	var batches [][]K
	for i := 0; i < len(slice); i += batch {
		end := i + batch

		if end > len(slice) {
			end = len(slice)
		}

		batches = append(batches, slice[i:end])
	}

	return batches
}

func PickField[K interface{}, V comparable](iterator []K, returner func(K) V) (result []V) {
	for _, item := range iterator {
		result = append(result, returner(item))
	}
	return
}

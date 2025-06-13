package utils

func UniqueSliceUInts(input []uint) []uint {
	seen := make(map[uint]bool)
	result := make([]uint, 0)

	for _, val := range input {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}

	return result
}

package typex

func IndexInt(items []int, v int) int {
	if len(items) == 0 {
		return -1
	}

	for index, item := range items {
		if item == v {
			return index
		}
	}

	return -1
}

func ContainInt(items []int, v int) bool {
	return IndexInt(items, v) != -1
}

func IntSliceRange(start, end, step int) []int {
	if step <= 0 {
		step = 1
	}

	slice := make([]int, (end-start+step)/step)
	index := 0
	for i := start; i <= end; i += step {
		slice[index] = i
		index++
	}

	return slice
}

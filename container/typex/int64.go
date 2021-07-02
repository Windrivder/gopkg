package typex

func IndexInt64(items []int64, v int64) int {
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

func ContainInt64(items []int64, v int64) bool {
	return IndexInt64(items, v) != -1
}

func NewInt64SliceRange(start, end, step int) []int64 {
	if step <= 0 {
		step = 1
	}

	slice := make([]int64, (end-start+step)/step)
	index := 0
	for i := start; i <= end; i += step {
		slice[index] = int64(i)
		index++
	}

	return slice
}

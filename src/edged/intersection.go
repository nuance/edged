package main

type sortedIndex []int64

func makeSortedIndex() sortedIndex {
	return []int64{}
}

func (si *sortedIndex) Add(id int64) {
	*si = append(*si, id)
}

func (a sortedIndex) intersect(b sortedIndex) sortedIndex {
	if len(a) == 0 || len(b) == 0 {
		return makeSortedIndex()
	}

	result := []int64{}
	for {
		if a[0] == b[0] {
			result = append(result, a[0])

			a = a[1:]
			b = b[1:]
		} else if a[0] < b[0] {
			a = a[1:]
		} else {
			b = b[1:]
		}

		if len(a) == 0 || len(b) == 0 {
			return sortedIndex(result)
		}
	}

	return makeSortedIndex()
}

package main

type sortedIndex []id

func (si *sortedIndex) add(id id) {
	*si = append(*si, id)
}

func (a sortedIndex) intersect(b sortedIndex) sortedIndex {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}

	result := []id{}
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

	return nil
}

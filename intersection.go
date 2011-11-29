package main

func intersect(a, b []int64) []int64 {
	if len(a) == 0 || len(b) == 0 {
		return empty
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
			return result
		}
	}

	return empty
}

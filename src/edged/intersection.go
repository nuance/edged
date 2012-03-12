package main

import (
	"sort"
)

func intersect(a, b []int64) []int64 {
	if len(a) == 0 || len(b) == 0 {
		return []int64{}
	}

	result := []int64{}
	for {
		if a[0] == b[0] {
			result = append(result, a[0])

			a = a[1:]
			b = b[1:]
		} else if a[0] < b[0] {
			next := sort.Search(len(a), func(i int) bool { return a[i] > b[0] })
			a = a[next:]
		} else {
			next := sort.Search(len(b), func(i int) bool { return b[i] > a[0] })
			b = b[next:]
		}

		if len(a) == 0 || len(b) == 0 {
			return result
		}
	}

	return []int64{}
}

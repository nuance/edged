package main

var empty = []int64{}

type PairIndex map[string] map[string] []int64

func (pi PairIndex) Get(a, b string) []int64 {
	if _, ok := pi[a]; !ok {
		return empty
	} else if _, ok := pi[a][b]; !ok {
		return empty
	}

	return pi[a][b]
}

func (pi PairIndex) Add(a, b string, id int64) {
	if _, ok := pi[a]; !ok {
		pi[a] = map[string] []int64{}
	}

	if _, ok := pi[b]; !ok {
		pi[b] = map[string] []int64{}
	}

	if _, ok := pi[a][b]; !ok {
		pi[a][b] = []int64{}
	}

	pi[a][b] = append(pi[a][b], id)
	pi[b][a] = pi[a][b]
}

func (pi PairIndex) Update(a, b string, id int64) {
	_, hasA := pi[a]
	_, hasB := pi[b]

	if !hasA && !hasB {
		return
	}

	pi.Add(a, b, id)
}

type IndexSet struct {
	// contains key(el, val) => doc ids
	indexes map[string] []int64
	// contains key(el, val) => key(other_el, other_val) => doc_ids
	intersections PairIndex
}

func EmptyIndexSet() *IndexSet {
	return &IndexSet{indexes: map[string] []int64{}, intersections: map[string] map[string] []int64{}}
}

func (is IndexSet) Lookup(key string) []int64 {
	if idx, ok := is.indexes[key]; ok {
		return idx
	}

	return empty
}

func (is IndexSet) Intersection(a, b string) []int64 {
	// This will only exist if one is a vip
	isect := is.intersections.Get(a, b)
	if len(isect) > 0 {
		return isect
	}

	// neither is a vip, so compute the intersection on the fly
	return intersect(is.Lookup(a), is.Lookup(b))
}

const IMPORTANT = 30

func (is *IndexSet) Add(node Node) {
	tokens := node.Tokens()

	// possibly create intersection indexes
	for _, token := range tokens {
		if len(is.indexes[token]) == IMPORTANT {
			for other, _ := range is.indexes {
				if token == other {
					continue
				}

				if len(is.intersections.Get(token, other)) == 0 {
					continue
				}

				for _, docId := range is.Intersection(token, other) {
					is.intersections.Add(token, other, docId)
				}
			}
		}
	}

	// add the token to the standard indexes. Do this after you create new
	// intersection indexes so that the next update step is consistent.
	for _, token := range tokens {
		if _, ok := is.indexes[token]; !ok {
			is.indexes[token] = []int64{}
		}
		is.indexes[token] = append(is.indexes[token], *node.Id)
	}

	// update intersection indexes with pairs
	for idx, token := range tokens {
		for _, other := range tokens[idx+1:] {
			is.intersections.Update(token, other, *node.Id)
		}
	}
}

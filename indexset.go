package main

var empty = []int64{}

type piKey struct {
	a, b string
}

func key(a, b string) piKey {
	if a >= b {
		return piKey{a, b}
	}
	return piKey{b, a}
}

type PairIndex map[piKey] []int64

func (pi PairIndex) Contains(a, b string) bool {
	_, ok := pi[key(a, b)]

	return ok
}

func (pi PairIndex) Get(a, b string) []int64 {
	val, ok := pi[key(a, b)]
	if !ok {
		return empty
	}

	return val
}

func (pi PairIndex) Add(a, b string, id int64) {
	k := key(a, b)
	pi[k] = append(pi[k], id)
}

func (pi PairIndex) Set(a, b string, ids []int64) {
	k := key(a, b)
	if _, ok := pi[k]; ok {
		panic("key already exists")
	}

	pi[k] = ids
}

type IndexSet struct {
	// contains key(el, val) => doc ids
	indexes map[string] []int64
	// contains key(el, val) => key(other_el, other_val) => doc_ids
	intersections PairIndex
}

func EmptyIndexSet() *IndexSet {
	return &IndexSet{indexes: map[string] []int64{}, intersections: PairIndex{}}
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

				if is.intersections.Contains(token, other) {
					continue
				}

				is.intersections.Set(token, other, is.Intersection(token, other))
			}
		}
	}

	// add the token to the standard indexes. Do this after you create new
	// intersection indexes so that the next update step is consistent.
	for _, token := range tokens {
		if _, ok := is.indexes[token]; !ok {
			is.indexes[token] = []int64{}
		}
		is.indexes[token] = append(is.indexes[token], node.Id)
	}

	// update intersection indexes with pairs
	for idx, token := range tokens {
		for _, other := range tokens[idx+1:] {
			if is.intersections.Contains(token, other) {
				is.intersections.Add(token, other, node.Id)
			}
		}
	}
}

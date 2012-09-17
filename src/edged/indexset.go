package main

type TokenPair [2]Token

func MakeTokenPair(a, b Token) TokenPair {
	if a.el > b.el || (a.el == b.el && a.id >= b.id) {
		return TokenPair{a, b}
	}
	return TokenPair{b, a}
}

type PairIndex map[TokenPair]sortedIndex

func (pi PairIndex) Contains(a, b Token) bool {
	_, ok := pi[MakeTokenPair(a, b)]

	return ok
}

func (pi PairIndex) Get(a, b Token) ([]int64, bool) {
	val, ok := pi[MakeTokenPair(a, b)]
	return val, ok
}

func (pi PairIndex) Add(a, b Token, id int64) {
	k := MakeTokenPair(a, b)
	idx := pi[k]
	idx.Add(id)
	pi[k] = idx
}

func (pi PairIndex) Set(a, b Token, ids []int64) {
	k := MakeTokenPair(a, b)
	if _, ok := pi[k]; ok {
		panic("key already exists")
	}

	pi[k] = ids
}

type IndexSet struct {
	id    map[Token]sortedIndex
	value map[string]int64
	vips  []Token
	// contains key(el, val) => key(other_el, other_val) => doc_ids
	intersection PairIndex
}

func EmptyIndexSet() *IndexSet {
	i := IndexSet{}

	i.id = map[Token]sortedIndex{}
	i.value = map[string]int64{}
	i.vips = []Token{}
	i.intersection = PairIndex{}

	return &i
}

func (is IndexSet) LookupValue(val string) (int64, bool) {
	id, ok := is.value[val]
	return id, ok
}

func (is IndexSet) LookupToken(a Token) (sortedIndex, bool) {
	ids, ok := is.id[a]
	return ids, ok
}

// Compute the intersection for two tokens. Doesn't acquire any locks.
func (is IndexSet) computeIntersection(a, b Token) []int64 {
	return is.id[a].intersect(is.id[b])
}

func (is *IndexSet) IntersectTokens(tokens []Token) []int64 {
	if len(tokens) == 0 {
		return []int64{}
	}

	running, ok := is.LookupToken(tokens[0])
	if !ok {
		return []int64{}
	}

	if len(tokens) == 1 {
		return running
	}

	for _, token := range tokens[1:] {
		next, ok := is.LookupToken(token)
		if !ok {
			return []int64{}
		}

		running = running.intersect(next)
	}

	return running
}

const IMPORTANT = 30

func (is *IndexSet) makeVip(token Token) {
	for _, other := range is.vips {
		if is.intersection.Contains(token, other) {
			return
		}

		is.intersection.Set(token, other, is.id[token].intersect(is.id[other]))
	}

	is.vips = append(is.vips, token)
}

func (is *IndexSet) Add(node Node) {
	if _, ok := is.value[node.Value]; ok {
		panic("Value already set")
	}
	is.value[node.Value] = node.Id

	tokens := node.Tokens()

	// ensure id lists for each token
	for _, token := range tokens {
		if _, ok := is.id[token]; !ok {
			is.id[token] = makeSortedIndex()
		}
	}

	// possibly create intersection indexes
	for _, token := range tokens {
		if len(is.id[token]) == IMPORTANT {
			is.makeVip(token)
		}
	}

	// add the token to the standard indexes. Do this after you create new
	// intersection indexes so that the next update step is consistent.
	for _, token := range tokens {
		idx := is.id[token]
		idx.Add(node.Id)
		is.id[token] = idx
	}

	// update intersection indexes with pairs
	for idx, token := range tokens {
		for _, other := range tokens[idx+1:] {
			if is.intersection.Contains(token, other) {
				is.intersection.Add(token, other, node.Id)
			}
		}
	}
}

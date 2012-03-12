package main

import (
	"sync"
)

type TokenPair [2]Token

func MakeTokenPair(a, b Token) TokenPair {
	if a.el > b.el {
		return TokenPair{a, b}
	} else if a.el == b.el && a.id >= b.id {
		return TokenPair{a, b}
	}
	return TokenPair{b, a}
}

type PairIndex map[TokenPair][]int64

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
	pi[k] = append(pi[k], id)
}

func (pi PairIndex) Set(a, b Token, ids []int64) {
	k := MakeTokenPair(a, b)
	if _, ok := pi[k]; ok {
		panic("key already exists")
	}

	pi[k] = ids
}

type IndexSet struct {
	id map[Token][]int64
	value map[string]int64
	// contains key(el, val) => key(other_el, other_val) => doc_ids
	intersection PairIndex

	lock sync.RWMutex
}

func EmptyIndexSet() *IndexSet {
	i := IndexSet{}

	i.id = map[Token][]int64{}
	i.value = map[string]int64{}
	i.intersection = PairIndex{}

	return &i
}

func (is IndexSet) LookupValue(val string) (int64, bool) {
	id, ok := is.value[val]
	return id, ok
}

func (is IndexSet) LookupToken(a Token) ([]int64, bool) {
	ids, ok := is.id[a]
	return ids, ok
}

// Compute the intersection for two tokens. Doesn't acquire any locks.
func (is IndexSet) computeIntersection(a, b Token) []int64 {
	aIds, ok := is.LookupToken(a)
	if !ok {
		return []int64{}
	}

	bIds, ok := is.LookupToken(b)
	if !ok {
		return []int64{}
	}

	return intersect(aIds, bIds)
}

func (is *IndexSet) IntersectTokens(tokens []Token) []int64 {
	if len(tokens) == 0 {
		return []int64{}
	}

	is.lock.RLock()
	defer is.lock.RUnlock()

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

		running = intersect(running, next)
	}

	return running
}

const IMPORTANT = 30

func (is *IndexSet) Add(node Node) {
	is.lock.Lock()
	defer is.lock.Unlock()

	is.value[node.Value] = node.Id

	tokens := node.Tokens()

	// ensure id lists for each token
	for _, token := range tokens {
		if _, ok := is.id[token]; !ok {
			is.id[token] = []int64{}
		}
	}

	// possibly create intersection indexes
	for _, token := range tokens {
		if len(is.id[token]) == IMPORTANT {
			for other, _ := range is.id {
				if token == other {
					continue
				}

				if is.intersection.Contains(token, other) {
					continue
				}

				is.intersection.Set(token, other, is.computeIntersection(token, other))
			}
		}
	}

	// add the token to the standard indexes. Do this after you create new
	// intersection indexes so that the next update step is consistent.
	for _, token := range tokens {
		is.id[token] = append(is.id[token], node.Id)
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

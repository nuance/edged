package main

import (
	"fmt"
)

type tokenPair [2]Token

func makeTokenPair(a, b Token) tokenPair {
	if a > b {
		return tokenPair{a, b}
	}
	return tokenPair{b, a}
}

type elIndex [1 << elBits]*sortedIndex

type indexSet struct {
	// Value => id hash index
	value map[string]id

	// id->el->idx lookup indexes
	token []elIndex

	// Materialized intersection indexes
	vips map[Token]bool
	pair map[tokenPair]*sortedIndex
}

func makeIndexSet() *indexSet {
	is := indexSet{}

	is.value = map[string]id{}

	// include an empty elIndex for 0 to simplify indexing (since ids start at 1)
	is.token = []elIndex{elIndex{}}

	is.vips = map[Token]bool{}
	is.pair = map[tokenPair]*sortedIndex{}

	return &is
}

func (is indexSet) lookupByValue(val string) id {
	return is.value[val]
}

func (is *indexSet) setByValue(val string, id id) {
	if _, ok := is.value[val]; ok {
		panic("Value is non-unique")
	}
	is.value[val] = id
}

func (is indexSet) lookupByToken(t Token) *sortedIndex {
	if int64(t.id()) < int64(len(is.token)) {
		return is.token[t.id()][t.el()]
	}
	return nil
}

func (is *indexSet) ensureByToken(t Token) {
	if int64(t.id()) == int64(len(is.token)) {
		is.token = append(is.token, elIndex{})
	} else if int64(t.id()) > int64(len(is.token)) {
		panic(fmt.Sprintf("id out-of-order: %d vs. current max %d", int64(t.id()), len(is.token)))
	}

	if is.token[t.id()][t.el()] == nil {
		is.token[t.id()][t.el()] = &sortedIndex{}
	}
}

// Compute the intersection for two tokens.
func (is indexSet) lookupPair(a, b Token) sortedIndex {
	if idx, ok := is.pair[makeTokenPair(a, b)]; ok {
		return *idx
	}

	aIdx := is.lookupByToken(a)
	bIdx := is.lookupByToken(b)
	if aIdx == nil || bIdx == nil {
		return nil
	}

	return aIdx.intersect(*bIdx)
}

func (is indexSet) lookupByTokens(tokens []Token) sortedIndex {
	if len(tokens) == 0 {
		return nil
	}

	first := is.lookupByToken(tokens[0])
	if first == nil {
		return nil
	} else if len(tokens) == 1 {
		return *first
	}

	running := *first
	for _, token := range tokens[1:] {
		next := is.lookupByToken(token)
		if next == nil {
			return nil
		}

		running = running.intersect(*next)
	}

	return running
}

const IMPORTANT = 30

func (is *indexSet) makeVip(token Token) {
	t := is.lookupByToken(token)
	for other, _ := range is.vips {
		si := t.intersect(*is.lookupByToken(other))
		is.pair[makeTokenPair(token, other)] = &si
	}

	is.vips[token] = true
}

func (is *indexSet) add(node Node) {
	if _, ok := is.value[node.Value]; ok {
		panic("Value already set")
	}
	is.value[node.Value] = node.Id

	tokens := node.Tokens()

	// ensure id lists for each token
	for _, token := range tokens {
		is.ensureByToken(token)
	}

	for _, token := range tokens {
		// possibly create intersection indexes
		if len(*is.lookupByToken(token)) == IMPORTANT {
			is.makeVip(token)
		}
	}

	// add the token to the standard indexes. Do this after you create new
	// intersection indexes so that the next update step is consistent.
	for _, token := range tokens {
		is.lookupByToken(token).add(node.Id)
	}

	// update intersection indexes with pairs
	for idx, token := range tokens {
		if _, ok := is.vips[token]; ok {
			for _, other := range tokens[idx+1:] {
				if pair := is.pair[makeTokenPair(token, other)]; pair != nil {
					pair.add(node.Id)
				}
			}
		}
	}
}

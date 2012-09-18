// additions to the node type

package main

import (
	"code.google.com/p/goprotobuf/proto"
	"edged/api"
	"sort"
)

type el int8
type id uint64

const (
	LEFT el = iota
	PROP
	RIGHT
	ID
)

// Number of bits need to represent the element
const elBits = 2

type Node struct {
	Id    id     `json:"id,omitempty"`
	Value string `json:"value,omitempty"`
	Edge  Edge   `json:"edge,omitempty"`
}

func (n *Node) FromApi(an *api.Node) {
	n.Id = id(*an.Id)
	n.Value = *an.Value

	if an.Edge != nil {
		n.Edge.FromApi(an.Edge)
	}
}

func (n Node) Api() *api.Node {
	res := &api.Node{}
	res.Id = proto.Int64(int64(n.Id))
	res.Value = proto.String(n.Value)
	res.Edge = n.Edge.Api()

	return res
}

type Edge struct {
	Left  id `json:"left,omitempty"`
	Prop  id `json:"prop,omitempty"`
	Right id `json:"right,omitempty"`
}

func (e *Edge) FromApi(ae *api.Edge) {
	e.Left = id(*ae.Left)
	e.Prop = id(*ae.Prop)
	e.Right = id(*ae.Right)
}

func (e *Edge) Api() *api.Edge {
	if e == nil {
		return nil
	}

	edge := &api.Edge{}
	edge.Left = proto.Int64(int64(e.Left))
	edge.Prop = proto.Int64(int64(e.Prop))
	edge.Right = proto.Int64(int64(e.Right))

	return edge
}

type Token uint64

func makeToken(id id, el el) Token {
	if (id >> (64 - elBits)) != 0 {
		panic("ID too large")
	}

	return Token(uint64(id) | (uint64(el) << (64 - elBits)))
}

func (t Token) id() id {
	return id((uint64(t) << elBits) >> elBits)
}

func (t Token) el() el {
	return el(t >> (64 - elBits))
}

type tokenSort []Token

func (ts tokenSort) Len() int {
	return len(ts)
}

func (ts tokenSort) Less(i, j int) bool {
	return ts[i].id() < ts[j].id()
}

func (ts tokenSort) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

var ts sort.Interface = tokenSort{}

func (n Node) Tokens() []Token {
	result := []Token{makeToken(n.Id, ID)}

	if n.Edge.Left != 0 {
		result = append(result, makeToken(n.Edge.Left, LEFT),
			makeToken(n.Edge.Prop, PROP),
			makeToken(n.Edge.Right, RIGHT))
	}

	sort.Sort(tokenSort(result))

	return result
}

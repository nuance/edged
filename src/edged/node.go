// additions to the node type

package main

import (
	"code.google.com/p/goprotobuf/proto"
	"edged/api"
)

type el int8

const (
	LEFT el = iota
	PROP
	RIGHT
	VALUE
	ID
)

type Node struct {
	Id    int64  `json:"id,omitempty"`
	Value string `json:"value,omitempty"`
	Edge  *Edge  `json:"edge,omitempty"`
}

func (n *Node) FromApi(an *api.Node) {
	n.Id = *an.Id
	n.Value = *an.Value

	if an.Edge == nil {
		n.Edge = nil
	} else {
		n.Edge = &Edge{}
		n.Edge.FromApi(an.Edge)
	}
}

func (n Node) Api() *api.Node {
	res := &api.Node{}
	res.Id = proto.Int64(n.Id)
	res.Value = proto.String(n.Value)
	res.Edge = n.Edge.Api()

	return res
}

type Edge struct {
	Left  int64 `json:"left,omitempty"`
	Prop  int64 `json:"prop,omitempty"`
	Right int64 `json:"right,omitempty"`
}

func (e *Edge) FromApi(ae *api.Edge) {
	e.Left = *ae.Left
	e.Prop = *ae.Prop
	e.Right = *ae.Right
}

func (e *Edge) Api() *api.Edge {
	if e == nil {
		return nil
	}

	edge := &api.Edge{}
	edge.Left = proto.Int64(e.Left)
	edge.Prop = proto.Int64(e.Prop)
	edge.Right = proto.Int64(e.Right)

	return edge
}

type Token struct {
	id int64
	el el
}

func (n Node) Tokens() []Token {
	result := []Token{{n.Id, ID}}

	if n.Edge != nil {
		result = append(result, Token{n.Edge.Left, LEFT}, Token{n.Edge.Prop, PROP}, Token{n.Edge.Right, RIGHT})
	}

	return result
}

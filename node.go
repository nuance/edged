// additions to the node type

package main

import (
	"fmt"
	"edged/api"
	"code.google.com/p/goprotobuf/proto"
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
	Id               int64  `json:"id,omitempty"`
	Value            string `json:"value,omitempty"`
	Edge             *Edge   `json:"edge,omitempty"`
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
	Left             int64 `json:"left,omitempty"`
	Prop             int64 `json:"prop,omitempty"`
	Right            int64 `json:"right,omitempty"`
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

var propKey = []string{"L", "P", "R", "V", "I"}

func ValueKey(element el, value string) string {
	return propKey[element] + value
}

func Key(element el, id int64) string {
	return propKey[element] + fmt.Sprintf("%d", id)
}

func (n Node) Tokens() []string {
	result := []string{Key(ID, n.Id), ValueKey(VALUE, n.Value)}

	if n.Edge != nil {
		result = append(result, Key(LEFT, n.Edge.Left), Key(PROP, n.Edge.Prop), Key(RIGHT, n.Edge.Right))
	}

	return result
}

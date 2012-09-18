package main

import (
	"fmt"
	"os"
	"testing"
)

func TestVipCreation(t *testing.T) {
	os.Remove("bench.graph")
	g, _ := Open("bench.graph")

	red, _ := g.Add(Node{Value: "red"})
	colorIs, _ := g.Add(Node{Value: "color_is"})

	vips := []id{}
	for i := 0; i < IMPORTANT+1; i++ {
		node, _ := g.Add(Node{Value: fmt.Sprintf("test%d", i)})

		edge := Edge{}
		edge.Left = id(node)
		edge.Prop = id(colorIs)
		edge.Right = id(red)
		vip, _ := g.Add(Node{Value: fmt.Sprintf("test%d is red", i), Edge: edge})
		vips = append(vips, id(vip))
	}

	pair := makeTokenPair(makeToken(id(red), RIGHT), makeToken(id(colorIs), PROP))
	if _, ok := g.indexes.pair[pair]; !ok {
		t.Error("no vip index for 'color_is' 'red'")
	}

	idxVip := g.indexes.pair[pair]
	if len(*idxVip) == len(vips) {
		for i := 0; i < len(vips); i++ {
			if (*idxVip)[i] != vips[i] {
				t.Error("vip index is incorrect")
				break
			}
		}
	} else {
		t.Error("vip index is incorrect")
	}

	computed := g.indexes.lookupByTokens([]Token{makeToken(id(red), RIGHT), makeToken(id(colorIs), PROP)})
	if len(computed) == len(vips) {
		for i := 0; i < len(vips); i++ {
			if computed[i] != vips[i] {
				t.Error("vip index is incorrect")
				break
			}
		}
	} else {
		t.Error("vip index is incorrect")
	}
}

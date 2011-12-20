package main

import (
	"os"
	"testing"
)

func TestVipCreation(t *testing.T) {
	os.Remove("bench.graph")
	g, _ := Open("bench.graph")

	red, _ := g.Add(Node{Value: "red"})
	color, _ := g.Add(Node{Value: "color_is"})

	vips := []int64{}
	for i := 0; i < IMPORTANT*2; i++ {
		node, _ := g.Add(Node{Value: "test"})

		edge := &Node_Edge{}
		edge.Left = node
		edge.Prop = color
		edge.Right = red
		vip, _ := g.Add(Node{Value: "", Edge: edge})
		vips = append(vips, vip)
	}

	if !g.Indexes.intersections.Contains(Key(RIGHT, red), Key(PROP, color)) {
		t.Error("no vip index if 'color_is' 'red'")
	}

	idxVip := g.Indexes.intersections.Get(Key(RIGHT, red), Key(PROP, color))
	if len(idxVip) == len(vips) {
		for i := 0; i < len(vips); i++ {
			if idxVip[i] != vips[i] {
				break
			}
		}

		return
	}

	t.Error("vip index is incorrect")
}

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

	vips := []int64{}
	for i := 0; i < IMPORTANT+1; i++ {
		node, _ := g.Add(Node{Value: fmt.Sprintf("test%d", i)})

		edge := &Edge{}
		edge.Left = node
		edge.Prop = colorIs
		edge.Right = red
		vip, _ := g.Add(Node{Value: fmt.Sprintf("test%d is red", i), Edge: edge})
		vips = append(vips, vip)
	}

	if !g.Indexes.intersection.Contains(Token{red, RIGHT}, Token{colorIs, PROP}) {
		t.Error("no vip index for 'color_is' 'red'")
	}

	idxVip := g.Indexes.intersection.Get(Token{red, RIGHT}, Token{colorIs, PROP})
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

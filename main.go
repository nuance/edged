package main

import (
	"fmt"
	"goprotobuf.googlecode.com/hg/proto"
)

func main() {
	g, err := Open("test.graph")
	if err != nil {
		panic(err)
	}
	if _, err := g.Add(Node{Value: proto.String("test")}); err != nil {
		panic(err)
	}

	for _, node := range g.Nodes {
		fmt.Printf("%s\n", node.String())
	}
}

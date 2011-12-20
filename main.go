package main

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
)

func main() {
	g, err := Open("test.graph")
	if err != nil {
		panic(err)
	}
	if _, err := g.Add(Node{Value: proto.String("test"), Edge: nil}); err != nil {
		panic(err)
	}

	for _, node := range g.Nodes {
		fmt.Printf("%s\n", node.String())
	}
}

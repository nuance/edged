package main

import (
	"fmt"
)

func main() {
	g, err := Open("test.graph")
	if err != nil {
		panic(err)
	}
	if _, err := g.Add(Node{Value: "test", Edge: nil}); err != nil {
		panic(err)
	}

	for _, node := range g.Nodes {
		fmt.Println(node)
	}
}

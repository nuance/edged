package main

import (
	"edged/api"
	"io"
	"os"
)

type Graph struct {
	Nodes   []Node
	Indexes *IndexSet
	log     *io.Writer
}

func Open(path string) (*Graph, error) {
	var log *io.ReadWriter
	if path != "" {
		l, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}

		var rw io.ReadWriter = l
		log = &rw
	}

	g := &Graph{}
	g.Nodes = []Node{}
	g.Indexes = EmptyIndexSet()
	if log != nil {
		w := io.Writer(*log)
		g.log = &w

		r := api.NodeReader{}
		n := &Node{}
		for {
			node, err := r.Read(*log)
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}

			n.FromApi(node)
			g.Nodes = append(g.Nodes, *n)
			g.Indexes.Add(*n)
		}
	}

	return g, nil
}

func (g *Graph) LookupValue(data string) *Node {
	id, ok := g.Indexes.LookupValue(data)
	if !ok {
		return nil
	}

	return &g.Nodes[id]
}

func (g *Graph) Add(node Node) (int64, error) {
	if n := g.LookupValue(node.Value); n != nil {
		return n.Id, nil
	}

	node.Id = int64(len(g.Nodes))
	g.Nodes = append(g.Nodes, node)

	if g.log != nil {
		if err := node.Api().Write(*g.log); err != nil {
			panic(err)
		}
	}

	g.Indexes.Add(node)
	return node.Id, nil
}

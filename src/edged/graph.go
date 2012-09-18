package main

import (
	"edged/api"
	"io"
	"os"
)

type Graph struct {
	maxId   id
	indexes *indexSet
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
	g.indexes = makeIndexSet()
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
			g.indexes.add(*n)
		}
	}

	return g, nil
}

func (g *Graph) Add(node Node) (int64, error) {
	if id := g.indexes.lookupByValue(node.Value); id != 0 {
		return int64(id), nil
	}

	g.maxId += 1
	node.Id = g.maxId

	if g.log != nil {
		if err := node.Api().Write(*g.log); err != nil {
			panic(err)
		}
	}

	g.indexes.add(node)
	return int64(node.Id), nil
}

package main

import (
	"io"
	"graphd/api"
	"os"
	"sync"
)

type Graph struct {
	Nodes   []Node
	Indexes *IndexSet
	log     io.Writer

	appendLock sync.Locker
}

func Open(path string) (*Graph, error) {
	log, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	g := &Graph{}
	g.Nodes = []Node{}
	g.Indexes = EmptyIndexSet()
	g.appendLock = &sync.Mutex{}
	g.log = log

	r := api.NodeReader{}
	n := &Node{}
	for {
		node, err := r.Read(log)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		n.FromApi(node)
		g.Nodes = append(g.Nodes, *n)
		g.Indexes.Add(*n)
	}

	return g, nil
}

func (g *Graph) Add(node Node) (int64, error) {
	g.appendLock.Lock()
	defer g.appendLock.Unlock()

	node.Id = int64(len(g.Nodes))

	g.Nodes = append(g.Nodes, node)

	if err := node.Api().Write(g.log); err != nil {
		panic(err)
	}

	g.Indexes.Add(node)
	return node.Id, nil
}

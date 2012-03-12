package main

import (
	"edged/api"
	"io"
	"os"
	"sync"
)

type Graph struct {
	Nodes   []Node
	Indexes *IndexSet
	log     io.Writer

	appendLock sync.Locker
}

type fileStub int

func (_ fileStub) Write(b []byte) (int, error) {
	return len(b), nil
}

func (_ fileStub) Read(b []byte) (int, error) {
	return 0, io.EOF
}

func Open(path string) (*Graph, error) {
	var log io.ReadWriter
	if path != "" {
		var err error
		log, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
	} else {
		log = fileStub(0)
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

func (g *Graph) LookupValue(data string) *Node {
	id, ok := g.Indexes.LookupValue(data)
	if !ok {
		return nil
	}

	return &g.Nodes[id]
}

func (g *Graph) Add(node Node) (int64, error) {
	g.appendLock.Lock()
	defer g.appendLock.Unlock()

	if n := g.LookupValue(node.Value); n != nil {
		return n.Id, nil
	}

	node.Id = int64(len(g.Nodes))

	g.Nodes = append(g.Nodes, node)

	if err := node.Api().Write(g.log); err != nil {
		panic(err)
	}

	g.Indexes.Add(node)
	return node.Id, nil
}

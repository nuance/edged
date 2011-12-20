package main

import (
	"code.google.com/p/goprotobuf/proto"
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

func Open(path string) (*Graph, error) {
	log, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	g := &Graph{Nodes: []Node{}, Indexes: EmptyIndexSet(), log: log, appendLock: &sync.Mutex{}}

	r := &NodeReader{}
	for {
		node, err := r.Read(log)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		g.Nodes = append(g.Nodes, *node)
		g.Indexes.Add(*node)
	}

	return g, nil
}

func (g *Graph) Add(node Node) (int64, error) {
	g.appendLock.Lock()
	defer g.appendLock.Unlock()

	node.Id = proto.Int64(int64(len(g.Nodes)))

	g.Nodes = append(g.Nodes, node)

	if err := node.Write(g.log); err != nil {
		panic(err)
	}

	g.Indexes.Add(node)
	return *node.Id, nil
}

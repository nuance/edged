package main

import (
	"encoding/binary"
	"io"
	"os"
	"sync"
	"goprotobuf.googlecode.com/hg/proto"
)

type Graph struct {
	Nodes []Node
	Indexes *IndexSet
	log io.Writer

	appendLock sync.Locker
}

func Open(path string) (*Graph, os.Error) {
	log, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	g := &Graph{Nodes: []Node{}, Indexes: EmptyIndexSet(), log: log, appendLock: &sync.Mutex{}}

	for {
		len := int16(0)
		if err := binary.Read(log, binary.LittleEndian, &len); err != nil {
			if err == os.EOF {
				break
			}

			panic(err)
		}
		buf := make([]byte, len)
		if _, err := log.Read(buf); err != nil {
			panic(err)
		}

		node := Node{}
		proto.Unmarshal(buf, &node)

		g.Nodes = append(g.Nodes, node)
		g.Indexes.Add(node)
	}

	return g, nil
}

func (g *Graph) Add(node Node) (int64, os.Error) {
	g.appendLock.Lock()
	defer g.appendLock.Unlock()

	node.Id = proto.Int64(int64(len(g.Nodes)))
	g.Nodes = append(g.Nodes, node)

	if data, err := proto.Marshal(&node); err != nil {
		return 0, err
	} else {
		if err := binary.Write(g.log, binary.LittleEndian, int16(len(data))); err != nil {
			return 0, err
		}

		if _, err := g.log.Write(data); err != nil {
			// need to rollback
			panic(err)
		}

		g.Indexes.Add(node)

		return *node.Id, nil
	}

	return *node.Id, nil
}

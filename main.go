package main

import (
	"encoding/binary"
	"sync"
	"io"
	"os"
	"goprotobuf.googlecode.com/hg/proto"
)

type Graph struct {
	nodes []Node
	log io.Writer

	appendLock sync.Locker
}

func Open(path string) (*Graph, os.Error) {
	log, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	g := &Graph{nodes: []Node{}, log: log, appendLock: &sync.Mutex{}}

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

		g.nodes = append(g.nodes, node)
	}

	return g, nil
}

func (g *Graph) Add(node Node) os.Error {
	g.appendLock.Lock()
	defer g.appendLock.Unlock()

	node.Id = proto.Int64(int64(len(g.nodes)))
	g.nodes = append(g.nodes, node)

	if data, err := proto.Marshal(&node); err != nil {
		return err
	} else {
		if err := binary.Write(g.log, binary.LittleEndian, int16(len(data))); err != nil {
			return err
		}

		_, err := g.log.Write(data)
		return err
	}

	return nil
}

func main() {
	g, err := Open("test.graph")
	if err != nil {
		panic(err)
	}
	if err := g.Add(Node{Value: proto.String("test")}); err != nil {
		panic(err)
	}
}

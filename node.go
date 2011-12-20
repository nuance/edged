// additions to the node type

package main

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"io"
)

type el int16

const (
	LEFT el = iota
	PROP
	RIGHT
	VALUE
	ID
)

func (node *Node) Write(w io.Writer) error {
	buf, err := proto.Marshal(node)
	if err != nil {
		panic(err)
	}

	if err := binary.WriteVarint(w, int64(len(buf))); err != nil {
		return err
	}

	_, err = bytes.NewBuffer(buf).WriteTo(w)
	return err
}

type NodeReader []byte

type byteReader struct {
	io.Reader
}

func (b byteReader) ReadByte() (byte, error) {
	buf := []byte{0}
	_, err := io.Reader(b).Read(buf)

	return buf[0], err
}

func (nr *NodeReader) Read(r io.Reader) (*Node, error) {
	l, err := binary.ReadVarint(byteReader{r})
	if err != nil {
		return nil, err
	}
	length := int(l)

	if *nr == nil || cap(*nr) < length {
		*nr = NodeReader(make([]byte, length))
	}
	buf := (*nr)[:length]

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	node := &Node{}
	if err := proto.Unmarshal(buf, node); err != nil {
		return nil, err
	}

	return node, nil
}

var propKey = []string{"L", "P", "R", "V", "I"}

func ValueKey(element el, value string) string {
	return propKey[element] + value
}

func Key(element el, id int64) string {
	return propKey[element] + string(id)
}

func (n Node) Tokens() []string {
	result := []string{Key(ID, *n.Id), ValueKey(VALUE, *n.Value)}

	if n.Edge != nil {
		result = append(result, Key(LEFT, *n.Edge.Left), Key(PROP, *n.Edge.Prop), Key(RIGHT, *n.Edge.Right))
	}

	return result
}

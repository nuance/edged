package api

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"io"
)

func (node *Node) Write(w io.Writer) error {
	buf, err := proto.Marshal(node)
	if err != nil {
		panic(err)
	}

	lenBuf := make([]byte, binary.MaxVarintLen64)
	l := binary.PutVarint(lenBuf, int64(len(buf)))
	if _, err := bytes.NewBuffer(lenBuf[:l]).WriteTo(w); err != nil {
		return err
	}

	_, err = bytes.NewBuffer(buf).WriteTo(w)
	return err
}

type byteReader struct {
	io.Reader
}

func (b byteReader) ReadByte() (byte, error) {
	buf := []byte{0}
	_, err := io.Reader(b).Read(buf)

	return buf[0], err
}

type NodeReader []byte

func (nr *NodeReader) Read(r io.Reader) (*Node, error) {
	l, err := binary.ReadVarint(byteReader{r})
	if err != nil {
		return nil, err
	}
	length := int(l)

	if cap(*nr) < length {
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

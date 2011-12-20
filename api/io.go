package api

import (
	"bytes"
	"encoding/binary"
	"io"
	"code.google.com/p/goprotobuf/proto"
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
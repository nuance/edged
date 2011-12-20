include $(GOROOT)/src/Make.inc

TARG=graphd

GOFILES=\
	graph.go \
	indexset.go \
	intersection.go \
	main.go \
	node.go \
	node.pb.go

include $(GOROOT)/src/Make.cmd
include $(GOROOT)/src/code.google.com/p/goprotobuf/Make.protobuf

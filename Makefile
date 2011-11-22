include $(GOROOT)/src/Make.inc

TARG=graphd

GOFILES=\
	main.go \
	node.pb.go

include $(GOROOT)/src/Make.cmd
include $(GOROOT)/src/pkg/goprotobuf.googlecode.com/hg/Make.protobuf
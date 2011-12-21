include $(GOROOT)/src/Make.inc

TARG=graphd

GOFILES=\
	graph.go\
	indexset.go\
	intersection.go\
	main.go\
	node.go\
	freebase.go\

include $(GOROOT)/src/Make.cmd

package main

import (
	"flag"
)

var freebaseDumpPath = flag.String("freebase", "", "path to a freebase dump")
var logPath = flag.String("insert-log", "", "path to write inserts to")

func main() {
	flag.Parse()

	g, err := Open(*logPath)
	if err != nil {
		panic(err)
	}

	if *freebaseDumpPath != "" {
		if err := g.ReadFreebase(*freebaseDumpPath); err != nil {
			panic(err)
		}
	}
}

package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var freebaseDumpPath = flag.String("freebase", "", "path to a freebase dump")
var logPath = flag.String("insert-log", "", "path to write inserts to")
var cpuProfile = flag.String("cpu", "", "write cpu profile to file")
var heapProfile = flag.String("heap", "", "write heap profile to file")

func main() {
	flag.Parse()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	g, err := Open(*logPath)
	if err != nil {
		panic(err)
	}

	if *freebaseDumpPath != "" {
		if err := g.ReadFreebase(*freebaseDumpPath); err != nil {
			panic(err)
		}
	}

	log.Println("done!")

	if *heapProfile != "" {
		f, err := os.Create(*heapProfile)
		if err != nil {
			log.Fatal(err)
		}
		runtime.GC()
		pprof.WriteHeapProfile(f)
	}
}

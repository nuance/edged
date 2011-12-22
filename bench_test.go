package main

import (
	"os"
	"testing"
)

func BenchmarkSerialInserts(b *testing.B) {
	os.Remove("bench.graph")
	g, _ := Open("bench.graph")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := g.Add(Node{Value: "test" + string(i)}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkOpenDB(b *testing.B) {
	os.Remove("bench.graph")
	g, _ := Open("bench.graph")

	for i := 0; i < b.N; i++ {
		if _, err := g.Add(Node{Value: "test" + string(i)}); err != nil {
			panic(err)
		}
	}

	b.ResetTimer()
	Open("bench.graph")
}

package main

import (
	"os"
	"testing"
)

func BenchmarkSerialInserts(b *testing.B) {
	os.Remove("bench.graph")
	g, err := Open("bench.graph")
	if err != nil {
		b.Fatalf("open failed: ", err.Error())
	}
	defer os.Remove("bench.graph")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := g.Add(Node{Value: "test" + string(i)}); err != nil {
			b.Fatalf(err.Error())
		}
	}
}

func BenchmarkOpenDB(b *testing.B) {
	os.Remove("bench.graph")
	g, err := Open("bench.graph")
	if err != nil {
		b.Fatalf("open failed: ", err.Error())
	}
	defer os.Remove("bench.graph")

	for i := 0; i < b.N; i++ {
		if _, err := g.Add(Node{Value: "test" + string(i)}); err != nil {
			b.Fatalf(err.Error())
		}
	}

	b.ResetTimer()
	Open("bench.graph")
}

package main

import "testing"

func Benchmark_EvenByModule(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvenByModulo(i)
	}
}

func Benchmark_EvenByBi(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvenByBit(i)
	}
}

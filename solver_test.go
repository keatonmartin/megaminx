package main

import (
	"testing"
)

// BenchmarkSolve benchmarks the Solve function
func BenchmarkSolve(b *testing.B) {
	// Create a sample Megaminx state with k=10 (adjust as needed)

	b.ResetTimer() // Reset the timer before starting the benchmark
	for i := 0; i < b.N; i++ {
		Test(13)
	}
}

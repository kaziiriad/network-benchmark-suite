// Command bench runs network benchmarks and reports results.
//
// Phase 1 hello-world: host mode, iperf3, throughput only.
// Run with: sudo ./bench -iter 5
package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/kaziiriad/network-benchmark-suite/internal/benchmark"
)

func main() {
	// flag.Int creates an -iter flag with default 5 and a help string.
	// flag.Parse() reads os.Args[1:]. Anything after the first non-flag
	// is a positional arg (we don't use any).
	iter := flag.Int("iter", 5, "number of iterations to run")
	flag.Parse()

	fmt.Printf("running host mode benchmark, %d iterations\n", *iter)

	// Collect throughput samples. Pre-allocated capacity (make([]T, 0, n))
	// avoids reallocation as we append — small thing, but it's the
	// idiomatic Go pattern for "I know exactly how many elements."
	results := make([]float64, 0, *iter)
	for i := 0; i < *iter; i++ {
		// RunIperf3 takes (host, durationSec). localhost is fine for
		// host mode — there's no namespace, the server is just the
		// local kernel. 3-second test is short enough to be quick,
		// long enough to average out sub-second jitter.
		r, err := benchmark.RunIperf3("localhost", 3)
		if err != nil {
			log.Fatalf("iteration %d failed: %v", i+1, err)
		}
		fmt.Printf("  iter %d: %.2f Gbps\n", i+1, r.Throughput)
		results = append(results, r.Throughput)
	}

	median := medianOf(results)
	fmt.Printf("\nhost mode, %d iterations, median throughput: %.2f Gbps\n", *iter, median)
}

// medianOf returns the median of a float64 slice.
//
// Implementation note: we sort a COPY of the input, so the caller's
// slice is unchanged. For a tiny slice (5-10 elements) the copy is
// free; for huge slices you'd use a selection algorithm. We don't
// need that here.
//
// We use the lower-median for even-length slices (the (n/2)th
// element after sorting). For benchmarks this is fine — picking the
// upper-median or averaging the two middle values changes the result
// by less than the noise floor.
func medianOf(xs []float64) float64 {
	cp := make([]float64, len(xs))
	copy(cp, xs)
	sort.Float64s(cp)

	n := len(cp)
	if n%2 == 1 {
		return cp[n/2]
	}
	return (cp[n/2-1] + cp[n/2]) / 2
}

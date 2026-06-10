// Package benchmark runs network benchmark tools (iperf3, netperf, qperf)
// and returns structured results.
//
// Result is intentionally a single struct with all metrics, even though
// some tools don't measure all of them. The homogeneous shape lets the
// 6-mode comparison report work without per-tool normalization. Zero
// values in a field mean "this tool did not measure that metric."
package benchmark

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Result is the structured output of a single benchmark run.
//
// All metrics are at the parser-boundary units: Gbps for throughput,
// microseconds for latency and jitter. Callers never do unit math.
type Result struct {
	Throughput float64 // Gbps
	Latency    float64 // microseconds (μs); 0 if not measured
	Jitter     float64 // microseconds (μs); 0 if not measured
}

// RunIperf3 runs `iperf3 -c host -J -t duration` and returns the
// throughput in Gbps.
//
// host is the iperf3 server to connect to. For Phase 1 hello-world
// this is always "localhost" — we'll add real servers in different
// network namespaces during Phase 2.
//
// durationSec is the test length in seconds. iperf3 picks default
// values for everything else (single TCP stream, default window).
func RunIperf3(host string, durationSec int) (Result, error) {
	// Build the command. We pass arguments as a slice (not a single
	// string) so the OS executes exactly these argv values — no
	// shell parsing, no injection risk, no quoting bugs.
	cmd := exec.Command(
		"iperf3",
		"-c", host, // client mode, target host
		"-J",      // JSON output (machine-readable)
		"-t", fmt.Sprint(durationSec), // test duration
	)

	// cmd.Output() runs the command, waits for it to finish, and
	// returns stdout as a byte slice. If the command exits non-zero,
	// err is *exec.ExitError and stderr is in err.Stderr.
	out, err := cmd.Output()
	if err != nil {
		return Result{}, fmt.Errorf("iperf3 failed: %w (stderr: %s)", err, stderrOf(err))
	}

	// Parse just the fields we need. Anonymous struct + json tags
	// is the Go idiom for "I want exactly these JSON keys, and I
	// don't need a named type for this." It mirrors iperf3's
	// actual JSON structure: { "end": { "sum_received": { ... } } }.
	// You can see the real shape by running `iperf3 -c localhost -J`.
	var parsed struct {
		End struct {
			SumReceived struct {
				BitsPerSecond float64 `json:"bits_per_second"`
			} `json:"sum_received"`
		} `json:"end"`
	}
	if err := json.Unmarshal(out, &parsed); err != nil {
		return Result{}, fmt.Errorf("parse iperf3 json: %w", err)
	}

	// Convert bits/sec → Gigabits/sec at the boundary. Callers
	// work in Gbps from here on.
	return Result{
		Throughput: parsed.End.SumReceived.BitsPerSecond / 1e9,
	}, nil
}

// stderrOf extracts stderr from an *exec.ExitError if possible.
// Other error types (binary not found, etc.) have no stderr to
// return, so we fall back to err.Error().
//
// We use type assertion (the `if ee, ok := ...` pattern) rather
// than a type switch because we only care about the one error
// type that actually carries stderr.
func stderrOf(err error) string {
	if ee, ok := err.(*exec.ExitError); ok {
		return string(ee.Stderr)
	}
	return err.Error()
}

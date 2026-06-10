# Network Performance Benchmarking Suite — Implementation Plan

## Context

Sultan Mahmud is a DevOps/Platform engineer building portfolio projects that demonstrate production-grade thinking. This project (#7 from the AWS DevOps Project Ideas catalogue) is a **learning project** covering three domains simultaneously: Go systems programming, Linux networking internals, and performance benchmarking methodology.

Target environment: **bare metal Linux** (laptop/desktop, root access, 2+ cores). This gives clean benchmark results and full access to network namespaces, TC, and all virtualization modes.

---

## Phase 1: High-Level Architecture

### What We're Building

A CLI tool that:
1. Creates network namespaces with different virtual networking modes (host, bridge, macvlan, ipvlan, veth+bridge, VXLAN)
2. Runs throughput + latency benchmarks inside each namespace using iperf3, netperf, qperf
3. Applies tc netem impairment rules and re-runs
4. Generates Markdown comparison tables + matplotlib charts

### Component Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         bench (Go CLI)                          │
├─────────────────────────────────────────────────────────────────┤
│  cmd/bench/                                                     │
│  ├── main.go           # Entry: flag parsing, orchestrate       │
│  ├── runner.go        # Concurrent goroutine-per-mode           │
│  ├── namespace.go      # ip netns + network mode setup/teardown  │
│  ├── benchmark.go      # iperf3/netperf/qperf invocation         │
│  ├── netem.go          # tc netem via netlink (no shell)         │
│  └── parser.go         # Parse tool output → structured types    │
├─────────────────────────────────────────────────────────────────┤
│  scripts/                                                       │
│  └── charts.py        # matplotlib charts (Python subprocess)    │
├─────────────────────────────────────────────────────────────────┤
│  reports/            # Markdown tables + PNG charts (output)     │
│  data/               # Raw JSON/CSV from benchmark tools         │
└─────────────────────────────────────────────────────────────────┘
```

### Why This Architecture

| Layer | Choice | Rationale |
|-------|--------|-----------|
| Orchestration | Go | Goroutine-per-mode parallelism, single binary, direct iperf3 library |
| Network setup | Go + netlink | No shell spawned, reliable teardown, educational |
| Benchmark tools | iperf3/netperf/qperf | Industry standard, cross-validate results |
| Chart generation | Python matplotlib | Industry standard, call from Go via exec |
| Output | Markdown + PNG | Human-readable + visual, ready for Medium article |

### Key Design Decisions

1. **Concurrent execution** — all 6 network modes run in parallel goroutines, each produces JSON results
2. **Structured result types** — Go structs for all metrics (throughput, latency, jitter), serialized to JSON
3. **Netem via netlink** — not shell `tc` — shows syscall-level understanding
4. **Python charts as subprocess** — Go writes CSV, Python reads CSV, produces PNG. Clean separation.
5. **No container orchestration** — direct `ip netns` commands, no Docker. This is bare-metal learning.

---

## Phase 2: Step-by-Step Deep Dive (TBD after Phase 1 approval)

Will be detailed per-module: file structure, function signatures, syscall approach, test strategy.

---

## Verification

1. On bare metal Linux, run `sudo ip netns add test` — verify namespace creation works
2. Run `which iperf3 netperf qperf` — verify benchmark tools are installable
3. Run `tc qdisc show` — verify TC netem is available
4. Execute full benchmark suite end-to-end, verify:
   - All 6 modes complete without error
   - Markdown report generates with all modes present
   - PNG charts render without Python errors

## Design Constraints

- **No Docker** — direct `ip netns` and `tc` syscalls for learning
- **Netem via netlink** — not shell `tc` — educational
- **Statistics** — median + stddev across iterations (not mean, outliers skew)
- **Hello world first** — Phase 1: single host mode end-to-end; Phase 2: 6-mode parallel
- **Concurrency** — goroutine-per-mode when scaling to 6 modes

## Learning Goals

1. Go systems programming (goroutines, netlink syscalls, CLI tools)
2. Linux networking internals (network namespaces, bridge/macvlan/ipvlan/veth/VXLAN, tc netem)
3. Performance benchmarking methodology (variable control, statistical rigor)

## Phase-by-Phase Learning Targets

### Phase 1 (host mode, single end-to-end)
By the end I should be able to explain:
- How Go invokes external processes and captures output
- Why we parse stdout as structured data instead of reading files
- What a network namespace is at the kernel level
- Why median + stddev is better than mean for benchmarks

### Phase 2 (6-mode parallel)
By the end I should be able to explain:
- How goroutines differ from threads
- Why we need synchronization when writing results concurrently
- What bridge vs macvlan vs ipvlan actually does to packets
- How tc netem intercepts traffic at the kernel level

## Decisions (Locked)

- **Approach:** Hello world first — single mode (host) end-to-end, then scale to 6-mode parallel
- **Statistics:** Median + Standard Deviation across multiple runs per mode
- **Report columns:** Mode | Throughput (Gbps) | Latency (μs) | Jitter (μs)

## Open Questions
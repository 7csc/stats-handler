# Go System Stats Collector

`stats` is a lightweight Go library that provides runtime stats, system information, and basic health check endpoints.

## Features

- Collect Go runtime information (Goroutine count, Memory stats, GC stats, etc.)
- Collect process information (PID, Executable path, etc.)
- Provide JSON stats via HTTP handler
- Health check handler
- Optional Prometheus exporter
- Configurable environment variable filtering for security

## Install

```sh
go get github.com/7csc/stats-handler
```

### Usage

```go
r := chi.NewRouter()
r.Get("/stats", stats.Handler)
r.Get("/healthz", stats.HealthzHandler)
```

### Example Output

```text
{
  "time": 1710000000000,
  "go_version": "go1.22",
  "go_os": "linux",
  "go_arch": "amd64",
  "cpus": 8,
  "goroutine_num": 15,
  "memory_alloc": 123456,
  "memory_total_alloc": 987654321,
  "memory_sys": 543210,
  "memory_usage_percent": 23.4,
  "gc_count": 12,
  "gc_pause_total_ns": 987654,
  "thread_count": 25,
  "pid": 12345,
  "ppid": 1,
  "executable": "/app/main",
  "uptime": 123456789,
  "env_vars": {
    "PATH": "/usr/bin:/bin",
    "SECRET_KEY": "******"
  }
}
```

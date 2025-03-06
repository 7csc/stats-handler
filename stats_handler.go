package stats

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

type Stats struct {
	Time             int64             `json:"time"`
	Version          string            `json:"go_version"`
	OS               string            `json:"go_os"`
	Arch             string            `json:"go_arch"`
	CPUs             int               `json:"cpus"`
	GoroutineNum     int               `json:"goroutine_num"`
	MemoryAlloc      uint64            `json:"memory_alloc"`
	MemoryTotalAlloc uint64            `json:"memory_total_alloc"`
	MemorySys        uint64            `json:"memory_sys"`
	MemoryUsage      float64           `json:"memory_usage_percent"`
	GCCount          uint32            `json:"gc_count"`
	GCPauseTotal     uint64            `json:"gc_pause_total_ns"`
	ThreadCount      int64             `json:"thread_count"`
	PID              int               `json:"pid"`
	PPID             int               `json:"ppid"`
	Executable       string            `json:"executable"`
	Uptime           int64             `json:"uptime"`
	EnvVars          map[string]string `json:"env_vars"`
	BuildInfo        string            `json:"build_info"`
}

var startTime int64
var envFilterFunc func(string, string) (string, string)

func init() {
	startTime = time.Now().Unix()
	SetEnvFilter(defaultEnvFilter)
}

func defaultEnvFilter(key, value string) (string, string) {
	if strings.Contains(strings.ToUpper(key), "SECRET") {
		return key, "*****"
	}
	return key, value
}

func SetEnvFilter(filter func(string, string) (string, string)) {
	envFilterFunc = filter
}

func CollectStats() *Stats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	totalMemory := mem.Sys
	usedMemory := mem.Alloc
	memoryUsage := (float64(usedMemory) / float64(totalMemory)) * 100

	envVars := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			key, value := envFilterFunc(pair[0], pair[1])
			envVars[key] = value
		}
	}

	execPath, _ := os.Executable()

	buildInfo := "unknown"
	if bi, ok := debug.ReadBuildInfo(); ok {
		buildInfo = bi.String()
	}

	return &Stats{
		Time:             time.Now().UnixNano(),
		Version:          runtime.Version(),
		OS:               runtime.GOOS,
		Arch:             runtime.GOARCH,
		CPUs:             runtime.NumCPU(),
		GoroutineNum:     runtime.NumGoroutine(),
		MemoryAlloc:      mem.Alloc,
		MemoryTotalAlloc: mem.TotalAlloc,
		MemorySys:        mem.Sys,
		MemoryUsage:      memoryUsage,
		GCCount:          mem.NumGC,
		GCPauseTotal:     mem.PauseTotalNs,
		ThreadCount:      runtime.NumCgoCall(),
		PID:              os.Getpid(),
		PPID:             os.Getppid(),
		Executable:       execPath,
		Uptime:           (time.Now().UnixNano() - startTime) / int64(time.Nanosecond),
		EnvVars:          envVars,
		BuildInfo:        buildInfo,
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	stats := CollectStats()

	w.Header().Set("Context-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

package stats

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"time"
)

type Stats struct {
	Time              int64             `json:"time"`
	Version           string            `json:"go_version"`
	OS                string            `json:"go_os"`
	Arch              string            `json:"go_arch"`
	CPUs              int               `json:"cpus"`
	GoroutineNum      int               `json:"goroutine_num"`
	MemoryAlloc       uint64            `json:"memory_alloc"`
	MemoryTotalAlloc  uint64            `json:"memory_total_alloc"`
	MemorySys         uint64            `json:"memory_sys"`
	MemoryUsage       float64           `json:"memory_usage_percent"`
	FileDescriptorNum int               `json:"file_descriptor_num"`
	EnvVars           map[string]string `json:"env_vars"`
	Uptime            int64             `json:"uptime"`
}

var startTime = time.Now().UnixNano()

func CollectStats() *Stats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	totalMemory := mem.Sys
	usedMemory := mem.Alloc
	memoryUsage := (float64(usedMemory) / float64(totalMemory)) * 100

	var rLimit syscall.Rlimit
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	fdUsage := int(rLimit.Cur) - int(rLimit.Max)

	envVars := make(map[string]string)
	for _, e := range os.Environ() {
		pair := splitEnvVar(e)
		envVars[pair[0]] = pair[1]
	}

	return &Stats{
		Time:              time.Now().UnixNano(),
		Version:           runtime.Version(),
		OS:                runtime.GOOS,
		Arch:              runtime.GOARCH,
		CPUs:              runtime.NumCPU(),
		GoroutineNum:      runtime.NumGoroutine(),
		MemoryAlloc:       mem.Alloc,
		MemoryTotalAlloc:  mem.TotalAlloc,
		MemorySys:         mem.Sys,
		MemoryUsage:       memoryUsage,
		FileDescriptorNum: fdUsage,
		EnvVars:           envVars,
		Uptime:            (time.Now().UnixNano() - startTime) / int64(time.Second),
	}
}

func splitEnvVar(env string) []string {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return []string{env[:i], env[i+1:]}
		}
	}
	return []string{env, ""}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	stats := CollectStats()

	w.Header().Set("Context-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

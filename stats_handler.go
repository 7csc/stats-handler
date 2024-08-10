package stats

import (
	"encoding/json"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Stats struct {
	Time             int64  `json:"time"`
	Version          string `json:"go_version"`
	OS               string `json:"go_os"`
	Arch             string `json:"go_arch"`
	CPUs             int    `json:"cpus"`
	GoroutineNum     int    `json:"goroutine_num"`
	MemoryAlloc      uint64 `json:"memory_alloc"`
	MemoryTotalAlloc uint64 `json:"memory_total_alloc"`
	MemorySys        uint64 `json:"memory_sys"`
}

var mux sync.Mutex

func CollectStats() *Stats {
	mux.Lock()
	defer mux.Unlock()

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	now := time.Now()

	return &Stats{
		Time:             now.UnixNano(),
		Version:          runtime.Version(),
		OS:               runtime.GOOS,
		Arch:             runtime.GOARCH,
		CPUs:             runtime.NumCPU(),
		GoroutineNum:     runtime.NumGoroutine(),
		MemoryAlloc:      mem.Alloc,
		MemoryTotalAlloc: mem.TotalAlloc,
		MemorySys:        mem.Sys,
	}

}

func Handler(w http.ResponseWriter, r *http.Request) {
	var jsonBytes []byte
	var jsonErr error

	var body string

	jsonBytes, jsonErr = json.Marshal(CollectStats())

	if jsonErr != nil {
		body = jsonErr.Error()
	} else {
		body = string(jsonBytes)
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Content-Length"] = strconv.Itoa(len(body))
	for name, val := range headers {
		w.Header().Set(name, val)
	}

	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	io.WriteString(w, body)
}

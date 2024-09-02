package stats_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/7csc/stats-handler"
)

func TestHandler(t *testing.T) {
	stats.StartTimeInit()
	os.Setenv("TEST_ENV", "test_value")
	defer os.Unsetenv("TEST_ENV")

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(stats.Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("API Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var responseStats stats.Stats
	err = json.Unmarshal(rr.Body.Bytes(), &responseStats)
	if err != nil {
		t.Errorf("Unable to parse response body: %v", err)
	}

	now := time.Now().UnixNano()
	if responseStats.Time > now || responseStats.Time < now-int64(time.Second) {
		t.Errorf("Unexpected time: got %v", responseStats.Time)
	}

	if responseStats.Version != runtime.Version() {
		t.Errorf("Unexpected Go version: got %v want %v", responseStats.Version, runtime.Version())
	}

	if responseStats.OS != runtime.GOOS {
		t.Errorf("Unexpected OS: got %v want %v", responseStats.OS, runtime.GOOS)
	}

	if responseStats.Arch != runtime.GOARCH {
		t.Errorf("Unexpected Arch: got %v want %v", responseStats.Arch, runtime.GOARCH)
	}

	if responseStats.CPUs != runtime.NumCPU() {
		t.Errorf("Unexpected CPUs: got %v want %v", responseStats.CPUs, runtime.NumCPU())
	}

	if responseStats.GoroutineNum != runtime.NumGoroutine() {
		t.Errorf("Unexpected GoroutineNum: got %v want %v", responseStats.GoroutineNum, runtime.NumGoroutine())
	}

	if responseStats.MemoryUsage < 0 || responseStats.MemoryUsage > 100 {
		t.Errorf("Unexpected MemoryUsage: got %v", responseStats.MemoryUsage)
	}

	if runtime.GOOS == "linux" {
		if responseStats.FileDescriptorNum <= 0 {
			t.Errorf("Unexpected FileDescriptorNum on Linux: got %v", responseStats.FileDescriptorNum)
		}
	} else {
		if responseStats.FileDescriptorNum != -1 {
			t.Errorf("Expected FileDescriptorNum to be -1 on non-Linux platforms, got %v", responseStats.FileDescriptorNum)
		}
	}

	if val, exists := responseStats.EnvVars["TEST_ENV"]; !exists || val != "test_value" {
		t.Errorf("Unexpected EnvVars: TEST_ENV not found or incorrect value, got %v", responseStats.EnvVars["TEST_ENV"])
	}

	if responseStats.Uptime <= 0 {
		t.Errorf("Unexpected Uptime: got %v", responseStats.Uptime)
	}

}

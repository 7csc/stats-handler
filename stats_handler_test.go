package stats_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/7csc/stats-handler"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()

	stats.Handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", res.StatusCode, http.StatusOK)
	}

	var data stats.Stats
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if data.CPUs <= 0 {
		t.Errorf("invalid CPU count: %d", data.CPUs)
	}

	if data.GoroutineNum <= 0 {
		t.Errorf("invalid goroutine count: %d", data.GoroutineNum)
	}

	if data.Uptime <= 0 {
		t.Errorf("invalid uptime: %d", data.Uptime)
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	stats.HealthHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", res.StatusCode, http.StatusOK)
	}

	body := w.Body.String()
	if strings.TrimSpace(body) != "OK" {
		t.Errorf("unexpected body: got %q, want %q", body, "OK")
	}
}

func TestEnvFiltert(t *testing.T) {
	stats.SetEnvFilter(func(key, value string) (string, string) {
		if strings.Contains(strings.ToUpper(key), "SECRET") {
			return key, "[FILTERD]"
		}
		return key, value
	})

	statsData := stats.CollectStats()

	for k, v := range statsData.EnvVars {
		if strings.Contains(strings.ToUpper(k), "SECRET") && v != "[FILTERD]" {
			t.Errorf("env var filtering failed: key=%s, value=%s", k, v)
		}
	}
}

func TestUptimeIncreases(t *testing.T) {
	initial := stats.CollectStats().Uptime
	time.Sleep(50 * time.Millisecond)
	later := stats.CollectStats().Uptime

	if later <= initial {
		t.Errorf("uptime didn't increase: start=%d, end=%d", initial, later)
	}
}

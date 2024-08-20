package stats

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("API Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var responseStats Stats
	err = json.Unmarshal(rr.Body.Bytes(), &responseStats)
	if err != nil {
		t.Errorf("Unable to parse response body: %v", err)
	}
}

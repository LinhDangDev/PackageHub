package auditor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckRuntimeEol_Node18_IsEol(t *testing.T) {
	// Mock endoflife.date API response for nodejs
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cycles := []EolCycle{
			{
				Cycle:   "22",
				Latest:  "22.14.0",
				EolDate: "2027-04-30",
			},
			{
				Cycle:   "18",
				Latest:  "18.20.8",
				EolDate: "2025-04-30", // In the past -> EOL
			},
		}
		_ = json.NewEncoder(w).Encode(cycles)
	}))
	defer mockServer.Close()

	client := NewAuditorClient(WithBaseUrl(mockServer.URL))

	item, err := client.CheckRuntimeEol("nodejs", "18.20.7")
	if err != nil {
		t.Fatalf("CheckRuntimeEol failed: %v", err)
	}

	if item.Status != "EOL_CRITICAL" {
		t.Errorf("expected status 'EOL_CRITICAL', got '%s'", item.Status)
	}
	if item.LatestVersion != "18.20.8" {
		t.Errorf("expected latest version '18.20.8', got '%s'", item.LatestVersion)
	}
}

func TestCheckRuntimeEol_Node22_IsActive(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cycles := []EolCycle{
			{
				Cycle:   "22",
				Latest:  "22.14.0",
				EolDate: "2027-04-30",
			},
		}
		_ = json.NewEncoder(w).Encode(cycles)
	}))
	defer mockServer.Close()

	client := NewAuditorClient(WithBaseUrl(mockServer.URL))

	item, err := client.CheckRuntimeEol("nodejs", "22.10.0")
	if err != nil {
		t.Fatalf("CheckRuntimeEol failed: %v", err)
	}

	if item.Status != "UPDATE_AVAILABLE" {
		t.Errorf("expected status 'UPDATE_AVAILABLE', got '%s'", item.Status)
	}
	if item.LatestVersion != "22.14.0" {
		t.Errorf("expected latest version '22.14.0', got '%s'", item.LatestVersion)
	}
}

func TestCheckNpmLatest(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"version": "2.1.261",
		})
	}))
	defer mockServer.Close()

	client := NewAuditorClient(WithNpmRegistryUrl(mockServer.URL))

	item, err := client.CheckNpmPackage("@anthropic-ai/claude-code", "2.1.246")
	if err != nil {
		t.Fatalf("CheckNpmPackage failed: %v", err)
	}

	if item.Status != "UPDATE_AVAILABLE" {
		t.Errorf("expected status 'UPDATE_AVAILABLE', got '%s'", item.Status)
	}
	if item.LatestVersion != "2.1.261" {
		t.Errorf("expected latest version '2.1.261', got '%s'", item.LatestVersion)
	}
}

func TestAuditorCache(t *testing.T) {
	cache := NewCache(1 * time.Minute)

	cache.Set("key1", "val1")
	val, ok := cache.Get("key1")
	if !ok || val != "val1" {
		t.Errorf("expected 'val1', got '%v'", val)
	}

	// Non-existent key
	_, ok = cache.Get("key2")
	if ok {
		t.Errorf("expected key2 to be missing")
	}
}

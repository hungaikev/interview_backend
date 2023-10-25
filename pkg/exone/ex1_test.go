package exone

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetCandidate(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		candidate := Candidate{ID: "123", FullName: "John Doe", Email: "john@example.com"}
		err := json.NewEncoder(w).Encode(candidate)
		if err != nil {
			return
		}
	}))
	defer server.Close()

	candidate, err := getCandidate(server.URL+"/candidates/", "123")
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	if candidate.ID != "123" || candidate.FullName != "John Doe" {
		t.Errorf("Unexpected candidate details: %+v", candidate)
	}
}

func TestGetJob(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		job := Job{ID: "456", Name: "Software Engineer", Description: "Develops software"}
		err := json.NewEncoder(w).Encode(job)
		if err != nil {
			return
		}
	}))
	defer server.Close()

	job, err := getJob(server.URL+"/jobs/", "456")
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	if job.ID != "456" || job.Name != "Software Engineer" {
		t.Errorf("Unexpected job details: %+v", job)
	}
}

func TestStoreCandidate(t *testing.T) {
	candidate := &Candidate{ID: "123", FullName: "John Doe", Email: "john@example.com"}
	err := storeCandidate(candidate)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	// Check if candidate was stored in mockDB
	storedCandidate, ok := mockDB["candidate_123"].(Candidate)
	if !ok || storedCandidate.FullName != "John Doe" {
		t.Errorf("Candidate was not stored correctly: %+v", storedCandidate)
	}
}

func TestStoreJob(t *testing.T) {
	job := &Job{ID: "456", Name: "Software Engineer", Description: "Develops software"}
	err := storeJob(job)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	// Check if job was stored in mockDB
	storedJob, ok := mockDB["job_456"].(Job)
	if !ok || storedJob.Name != "Software Engineer" {
		t.Errorf("Job was not stored correctly: %+v", storedJob)
	}
}

func TestWebhook(t *testing.T) {
	// Mock getCandidate and getJob API servers
	mockCandidateServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		candidate := Candidate{ID: "123", FullName: "John Doe", Email: "john@example.com"}
		err := json.NewEncoder(w).Encode(candidate)
		if err != nil {
			return
		}
	}))
	defer mockCandidateServer.Close()

	mockJobServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		job := Job{ID: "456", Name: "Software Engineer", Description: "Develops software"}
		err := json.NewEncoder(w).Encode(job)
		if err != nil {
			return
		}
	}))
	defer mockJobServer.Close()

	// Sample payload for the Webhook
	payload := HookPayload[Application]{
		Hook:          Hook{ID: "1", Resource: "test", Action: "update", Target: "sample"},
		LinkedAccount: "testAccount",
		Data:          Application{Candidate: "123", Job: "456"},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	// Invoke Webhook with mock server URLs
	if err := Webhook(c, mockCandidateServer.URL+"/candidates/", mockJobServer.URL+"/jobs/"); err != nil {
		t.Errorf("Webhook function failed with error %v", err)
	}

	// Check HTTP status code
	assert.Equal(t, http.StatusOK, rec.Code) // assuming that the Webhook function returns a 200 status on success
}

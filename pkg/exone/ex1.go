package exone

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

// We receive notification via a webhook, for candidates changes. The webhook contains only a little part of the resource we need in our platform.
// For each request we need to retrieve the relevant information needed and store them in a database.
// The operation of retrieval is syncronous.
// Lately we noticed an increased number of errors with http status 429.
// You are assigned to the team that needs to undestand and fix the problem

var mockDB = make(map[string]interface{})

type Hook struct {
	ID       string `json:"id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Target   string `json:"target"`
}

type HookPayload[T any] struct {
	Hook          Hook   `json:"hook"`
	LinkedAccount string `json:"linked_account"`
	Data          T      `json:"data"`
}

type BaseApplication[C string, J string, St string] struct {
	ID           string    `json:"id"`
	RemoteID     string    `json:"remote_id"`
	Candidate    C         `json:"candidate"`
	Job          J         `json:"job"`
	AppliedAt    time.Time `json:"applied_at"`
	RejectedAt   time.Time `json:"rejected_at"`
	Source       string    `json:"source"`
	CreditedTo   any       `json:"credited_to"`
	CurrentStage St        `json:"current_stage"`
	RejectReason any       `json:"reject_reason"`
	RemoteData   any       `json:"remote_data"`
}

type Application = BaseApplication[string, string, string]

type Candidate struct {
	ID       string
	FullName string
	Email    string
}

type Job struct {
	ID          string
	Name        string
	Description string
}

func Webhook(c echo.Context, candidateEndpoint, jobEndpoint string) error {
	var hp HookPayload[Application]
	if err := json.NewDecoder(c.Request().Body).Decode(&hp); err != nil {
		return err
	}

	// Create an errgroup for concurrent operations
	var g errgroup.Group
	var candidate *Candidate
	var job *Job

	// Concurrently retrieve the candidate information
	g.Go(func() error {
		var err error
		candidate, err = getCandidate(candidateEndpoint, hp.Data.Candidate)
		return err
	})

	// Concurrently retrieve the job information
	g.Go(func() error {
		var err error
		job, err = getJob(jobEndpoint, hp.Data.Job)
		return err
	})

	// Wait for both goroutines to finish and check for errors
	if err := g.Wait(); err != nil {
		return err
	}

	// Reset errgroup for storing data
	g = errgroup.Group{}

	// Concurrently store the candidate information
	g.Go(func() error {
		return storeCandidate(candidate)
	})

	// Concurrently store the job information
	g.Go(func() error {
		return storeJob(job)
	})

	// Wait for both goroutines to finish and check for errors
	return g.Wait()
}

func getCandidate(endpoint string, id string) (*Candidate, error) {
	resp, err := http.Get(endpoint + id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve candidate: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get candidate with status: %d", resp.StatusCode)
	}

	var candidate Candidate
	if err := json.NewDecoder(resp.Body).Decode(&candidate); err != nil {
		return nil, fmt.Errorf("failed to decode candidate: %w", err)
	}

	return &candidate, nil
}

func getJob(endpoint string, id string) (*Job, error) {
	resp, err := http.Get(endpoint + id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get job with status: %d", resp.StatusCode)
	}

	var job Job
	err = json.NewDecoder(resp.Body).Decode(&job)
	if err != nil {
		return nil, fmt.Errorf("failed to decode job response: %w", err)
	}

	return &job, nil
}

func storeCandidate(c *Candidate) error {
	// Simulate storing to a mock database
	key := "candidate_" + c.ID
	mockDB[key] = *c
	return nil
}

func storeJob(j *Job) error {
	// Simulate storing to a mock database
	key := "job_" + j.ID
	mockDB[key] = *j
	return nil
}

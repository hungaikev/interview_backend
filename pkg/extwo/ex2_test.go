package extwo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrape(t *testing.T) {
	// Mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><head><title>Mock Page</title></head><body>Mock body content<a href="http://example.com">Link</a></body></html>`))
	}))
	defer mockServer.Close()

	scrapper := NewCollyScrapper()
	data, err := scrapper.Scrape(mockServer.URL)

	assert.NoError(t, err)
	assert.Equal(t, "Mock Page", data.Title)
	assert.Equal(t, "Mock body contentLink", data.Body)
	assert.Equal(t, mockServer.URL, data.URL)
	assert.Contains(t, data.Links, "http://example.com")
}

func TestScrapeMultiple(t *testing.T) {
	// Create two mock HTTP servers
	mockServer1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><head><title>Mock Page 1</title></head><body>Content 1<a href="http://example1.com">Link1</a></body></html>`))
	}))
	defer mockServer1.Close()

	mockServer2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><head><title>Mock Page 2</title></head><body>Content 2<a href="http://example2.com">Link2</a></body></html>`))
	}))
	defer mockServer2.Close()

	scrapper := NewCollyScrapper()
	results, err := scrapper.ScrapeMultiple([]string{mockServer1.URL, mockServer2.URL})

	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Check data for mockServer1
	assert.Equal(t, "Mock Page 1", results[0].Title)
	assert.Equal(t, "Content 1Link1", results[0].Body)
	assert.Contains(t, results[0].Links, "http://example1.com")

	// Check data for mockServer2
	assert.Equal(t, "Mock Page 2", results[1].Title)
	assert.Equal(t, "Content 2Link2", results[1].Body)
	assert.Contains(t, results[1].Links, "http://example2.com")
}

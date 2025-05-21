package kbapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockTransport implements the Transport interface for testing
type MockTransport struct {
	// Mock method specifies what we expect to be returned
	MockResponse *http.Response
	MockError    error

	// RecordedRequests stores all requests for later verification
	RecordedRequests []*http.Request
}

// NewMockTransport creates a new MockTransport with the given response and error
func NewMockTransport(statusCode int, responseBody interface{}, err error) *MockTransport {
	var body []byte
	var respErr error

	if responseBody != nil {
		body, respErr = json.Marshal(responseBody)
		if respErr != nil {
			panic(respErr) // In tests, we can panic for unexpected marshaling errors
		}
	}

	return &MockTransport{
		MockResponse: &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		},
		MockError:        err,
		RecordedRequests: make([]*http.Request, 0),
	}
}

// NewMockTransportWithRawResponse creates a mock transport with a raw string response
func NewMockTransportWithRawResponse(statusCode int, responseBody string, err error) *MockTransport {
	return &MockTransport{
		MockResponse: &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(bytes.NewReader([]byte(responseBody))),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		},
		MockError:        err,
		RecordedRequests: make([]*http.Request, 0),
	}
}

// Perform implements the Transport interface
func (m *MockTransport) Perform(req *http.Request) (*http.Response, error) {
	// Record the request for later verification
	m.RecordedRequests = append(m.RecordedRequests, req)

	// Return the mocked response and error
	return m.MockResponse, m.MockError
}

// LastRequest returns the most recent request made through this transport
func (m *MockTransport) LastRequest() *http.Request {
	if len(m.RecordedRequests) == 0 {
		return nil
	}
	return m.RecordedRequests[len(m.RecordedRequests)-1]
}

// RequestCount returns the number of requests made
func (m *MockTransport) RequestCount() int {
	return len(m.RecordedRequests)
}

// AssertRequestPath asserts that the HTTP request path is as expected
func AssertRequestPath(t *testing.T, req *http.Request, expected string) {
	t.Helper()
	assert.Equal(t, expected, req.URL.Path, "Request path should match")
}

// AssertRequestMethod asserts that the HTTP request method is as expected
func AssertRequestMethod(t *testing.T, req *http.Request, expected string) {
	t.Helper()
	assert.Equal(t, expected, req.Method, "Request method should match")
}

// AssertRequestParam asserts that a query parameter exists with the expected value
func AssertRequestParam(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()
	values := req.URL.Query()[key]
	assert.Greater(t, len(values), 0, "Query param %s should exist", key)
	if len(values) > 0 {
		assert.Equal(t, expected, values[0], "Query param %s should match", key)
	}
}

// AssertRequestHeader asserts that a header exists with the expected value
func AssertRequestHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()
	values := req.Header.Values(key)
	assert.Greater(t, len(values), 0, "Header %s should exist", key)
	if len(values) > 0 {
		assert.Equal(t, expected, values[0], "Header %s should match", key)
	}
}

// AssertRequestBodyJSON asserts that the request body matches the expected JSON
func AssertRequestBodyJSON(t *testing.T, req *http.Request, expected interface{}) {
	t.Helper()

	// Read body
	bodyBytes, err := io.ReadAll(req.Body)
	assert.NoError(t, err, "Should be able to read request body")

	// Replace body for future reads
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// Parse JSON
	var actualJSON interface{}
	err = json.Unmarshal(bodyBytes, &actualJSON)
	assert.NoError(t, err, "Request body should be valid JSON")

	// Marshal expected to normalize it
	expectedBytes, err := json.Marshal(expected)
	assert.NoError(t, err, "Expected value should be marshalable")

	var expectedJSON interface{}
	err = json.Unmarshal(expectedBytes, &expectedJSON)
	assert.NoError(t, err, "Expected value should be valid JSON")

	// Compare
	assert.Equal(t, expectedJSON, actualJSON, "Request body should match expected")
}

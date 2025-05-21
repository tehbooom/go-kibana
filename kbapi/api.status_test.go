package kbapi

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_StatusGet(t *testing.T) {
	mockStatus := KibanaStatusResponse{}
	mockStatus.Status.Overall.Level = "available"
	mockStatus.Status.Overall.Summary = "Kibana is available"
	mockStatus.Version.Number = "8.0.0"
	mockStatus.Version.BuildHash = "abcdef123456"
	mockStatus.Metrics.CollectionIntervalInMillis = 5000
	mockStatus.Name = "test-kibana"
	mockStatus.Uuid = "9db06eba-a178-4de9-89de-cd5f08889810"

	testCases := []struct {
		name         string
		request      *GetStatusRequest
		statusCode   int
		responseBody interface{}
		responseErr  error
		queryParams  map[string]string
		customHeader http.Header
	}{
		{
			name:         "Success with no parameters",
			request:      &GetStatusRequest{},
			responseBody: mockStatus,
			statusCode:   200,
			responseErr:  nil,
		},
		{
			name:         "Success with header",
			request:      &GetStatusRequest{},
			responseBody: mockStatus,
			statusCode:   200,
			responseErr:  nil,
			customHeader: http.Header{
				"x-custom-header": []string{"custom-value"},
			},
		},
		{
			name:         "Success with parameters",
			request:      &GetStatusRequest{V7format: BoolPtr(false), V8format: BoolPtr(false)},
			responseBody: mockStatus,
			statusCode:   200,
			responseErr:  nil,
		},
		{
			name:    "Error with 500 status code as JSON",
			request: &GetStatusRequest{},
			responseBody: map[string]interface{}{
				"error":   "Internal Server Error",
				"message": "Something went wrong",
			},
			statusCode:  500,
			responseErr: nil,
		},
		{
			name:         "Error with 404 status code as string",
			request:      &GetStatusRequest{},
			responseBody: "Not Found", // Non-JSON response
			statusCode:   404,
			responseErr:  nil,
		},
		{
			name:         "Transport-level error",
			request:      &GetStatusRequest{},
			responseBody: nil,
			statusCode:   0,
			responseErr:  fmt.Errorf("Transport error: connection refused"),
		},
		{
			name:         "Success without request",
			request:      nil,
			responseBody: nil,
			statusCode:   500,
			responseErr:  fmt.Errorf("Error something went wrong"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			var mockTransport *MockTransport

			// For non-JSON error bodies, use the raw string response constructor
			if test.name == "Error with 404 status code as string" {
				mockTransport = NewMockTransportWithRawResponse(test.statusCode, test.responseBody.(string), test.responseErr)
			} else {
				mockTransport = NewMockTransport(test.statusCode, test.responseBody, test.responseErr)
			}

			api := &API{transport: mockTransport}

			var opts []RequestOption
			if test.customHeader != nil {
				opts = append(opts, WithHeaders(test.customHeader))
			}

			statusFunc := api.newStatusFunc()
			resp, err := statusFunc(context.Background(), test.request, opts...)

			// All error cases should return an error
			if test.responseErr != nil || test.statusCode > 200 {
				require.Error(t, err, "Expected an error for status code %d", test.statusCode)

				if test.responseErr != nil {
					// For transport errors, we should get the exact same error
					assert.Equal(t, test.responseErr, err)
				} else {
					// For HTTP status errors, check that the status code is in the error message
					errMsg := err.Error()
					assert.Contains(t, errMsg, fmt.Sprintf("HTTP Status Code %d", test.statusCode),
						"Error message should contain status code")

					// For JSON error responses, verify error object was set
					if test.name == "Error with 500 status code as JSON" {
						require.NotNil(t, resp)
						assert.NotNil(t, resp.Error, "Response error field should be set")
					}

					// For string error responses, verify error string was set
					if test.name == "Error with 404 status code as string" {
						require.NotNil(t, resp)
						assert.Equal(t, "Not Found", resp.Error)
					}
				}
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, test.statusCode, resp.StatusCode, "Expected status code %d, got %d", test.statusCode, resp.StatusCode)

			req := mockTransport.LastRequest()
			require.NotNil(t, req)
			AssertRequestMethod(t, req, http.MethodGet)
			AssertRequestPath(t, req, "/api/status")

			if test.customHeader != nil {
				for key, values := range test.customHeader {
					for _, value := range values {
						AssertRequestHeader(t, req, key, value)
					}
				}
			}

			for key, value := range test.queryParams {
				AssertRequestParam(t, req, key, value)
			}

			if resp != nil && resp.Body != nil {
				assert.Equal(t, mockStatus.Status.Overall.Level, resp.Body.Status.Overall.Level)
			}
		})
	}
}

func TestAPI_StatusGetRedacted(t *testing.T) {
	mockStatus := KibanaStatusRedactedResponse{}
	mockStatus.Status.Overall.Level = "available"

	testCases := []struct {
		name         string
		request      *GetStatusRequest
		statusCode   int
		responseBody interface{}
		responseErr  error
		queryParams  map[string]string
		customHeader http.Header
	}{
		{
			name:         "Success with no parameters",
			request:      &GetStatusRequest{},
			responseBody: mockStatus,
			statusCode:   200,
			responseErr:  nil,
		},
		{
			name:         "Success with header",
			request:      &GetStatusRequest{},
			responseBody: mockStatus,
			statusCode:   200,
			responseErr:  nil,
			customHeader: http.Header{
				"x-custom-header": []string{"custom-value"},
			},
		},
		{
			name:         "Success with parameters",
			request:      &GetStatusRequest{V7format: BoolPtr(false), V8format: BoolPtr(false)},
			responseBody: mockStatus,
			statusCode:   200,
			responseErr:  nil,
		},
		{
			name:    "Error with 500 status code as JSON",
			request: &GetStatusRequest{},
			responseBody: map[string]interface{}{
				"error":   "Internal Server Error",
				"message": "Something went wrong",
			},
			statusCode:  500,
			responseErr: nil,
		},
		{
			name:         "Error with 404 status code as string",
			request:      &GetStatusRequest{},
			responseBody: "Not Found", // Non-JSON response
			statusCode:   404,
			responseErr:  nil,
		},
		{
			name:         "Transport-level error",
			request:      &GetStatusRequest{},
			responseBody: nil,
			statusCode:   0,
			responseErr:  fmt.Errorf("Transport error: connection refused"),
		},
		{
			name:         "Success without request",
			request:      nil,
			responseBody: nil,
			statusCode:   500,
			responseErr:  fmt.Errorf("Error something went wrong"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			var mockTransport *MockTransport

			// For non-JSON error bodies, use the raw string response constructor
			if test.name == "Error with 404 status code as string" {
				mockTransport = NewMockTransportWithRawResponse(test.statusCode, test.responseBody.(string), test.responseErr)
			} else {
				mockTransport = NewMockTransport(test.statusCode, test.responseBody, test.responseErr)
			}

			api := &API{transport: mockTransport}

			var opts []RequestOption
			if test.customHeader != nil {
				opts = append(opts, WithHeaders(test.customHeader))
			}

			statusFunc := api.newStatusRedactedFunc()
			resp, err := statusFunc(context.Background(), test.request, opts...)

			// All error cases should return an error
			if test.responseErr != nil || test.statusCode > 200 {
				require.Error(t, err, "Expected an error for status code %d", test.statusCode)

				if test.responseErr != nil {
					// For transport errors, we should get the exact same error
					assert.Equal(t, test.responseErr, err)
				} else {
					// For HTTP status errors, check that the status code is in the error message
					errMsg := err.Error()
					assert.Contains(t, errMsg, fmt.Sprintf("HTTP Status Code %d", test.statusCode),
						"Error message should contain status code")

					// For JSON error responses, verify error object was set
					if test.name == "Error with 500 status code as JSON" {
						require.NotNil(t, resp)
						assert.NotNil(t, resp.Error, "Response error field should be set")
					}

					// For string error responses, verify error string was set
					if test.name == "Error with 404 status code as string" {
						require.NotNil(t, resp)
						assert.Equal(t, "Not Found", resp.Error)
					}
				}
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, test.statusCode, resp.StatusCode, "Expected status code %d, got %d", test.statusCode, resp.StatusCode)

			req := mockTransport.LastRequest()
			require.NotNil(t, req)
			AssertRequestMethod(t, req, http.MethodGet)
			AssertRequestPath(t, req, "/api/status")

			if test.customHeader != nil {
				for key, values := range test.customHeader {
					for _, value := range values {
						AssertRequestHeader(t, req, key, value)
					}
				}
			}

			for key, value := range test.queryParams {
				AssertRequestParam(t, req, key, value)
			}

			if resp != nil && resp.Body != nil {
				assert.Equal(t, mockStatus.Status.Overall.Level, resp.Body.Status.Overall.Level)
			}
		})
	}
}

package kbapi

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
)

// Instrumented allows to retrieve the current transport Instrumentation
type Instrumented elastictransport.Instrumented

// Instrumentation defines the interface for the instrumentation API.
type Instrumentation elastictransport.Instrumentation

// BoolPtr returns a pointer to v.
//
// It is used as a convenience function for converting a bool value
// into a pointer when passing the value to a function or struct field
// which expects a pointer.
func BoolPtr(v bool) *bool { return &v }

// IntPtr returns a pointer to v.
//
// It is used as a convenience function for converting an int value
// into a pointer when passing the value to a function or struct field
// which expects a pointer.
func IntPtr(v int) *int { return &v }

// Float32Ptr returns a pointer to v.
//
// It is used as a convenience function for converting an int value
// into a pointer when passing the value to a function or struct field
// which expects a pointer.
func Float32Ptr(v float32) *float32 { return &v }

// StrPtr returns a pointer to v.
//
// It is used as a convenience function for converting an string value
// into a pointer when passing the value to a function or struct field
// which expects a pointer.
func StrPtr(v string) *string { return &v }

// SliceStrPtr returns a pointer to v.
//
// It is used as a convenience function for converting an string value
// into a pointer when passing the value to a function or struct field
// which expects a pointer.
func SliceStrPtr(v []string) *[]string { return &v }

// PrettyPrint converts any response to pretty-formatted JSON
func PrettyPrint(response interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Error formatting response: %v\n", err)
	}
	return string(jsonData), nil
}

// PrintRawBody returns that rawBody response as a string
func PrintRawBody(rawBody io.ReadCloser) (string, error) {
	bodyBytes, err := io.ReadAll(rawBody)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

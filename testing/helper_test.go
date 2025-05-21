package testing

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/tehbooom/go-kibana"
)

// TestClient represents a test Kibana API client
type TestClient struct {
	KibanaClient     *kibana.Client
	ElasticsearchURL string
	KibanaURL        string
	Context          context.Context
	CancelFunc       context.CancelFunc
}

// NewTestClient creates a new test client
func NewTestClient(t *testing.T) *TestClient {
	t.Helper()

	// Get URLs from environment variables
	kibanaURL := os.Getenv(envKibanaURL)
	if kibanaURL == "" {
		kibanaURL = "https://localhost:5601"
	}

	esURL := os.Getenv(envElasticsearchURL)
	if esURL == "" {
		esURL = "https://localhost:9200"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Create the Kibana API client
	client, err := kibana.NewClient(kibana.Config{
		Addresses: []string{kibanaURL},
		Username:  "elastic",
		Password:  "changeme",
		Transport: transport,
	})
	if err != nil {
		t.Fatalf("Failed to create Kibana client: %v", err)
	}

	return &TestClient{
		KibanaClient:     client,
		ElasticsearchURL: esURL,
		KibanaURL:        kibanaURL,
		Context:          ctx,
		CancelFunc:       cancel,
	}
}

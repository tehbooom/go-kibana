package kibana

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/tehbooom/go-kibana/kbapi"
	"go.opentelemetry.io/otel/trace"
)

// Default values for client configuration
const (
	defaultURL = "http://localhost:5601"
	Version    = "0.1.0"
)

// Config represents the client configuration.
type Config struct {
	Addresses       []string    // A list of Kibana instances to use.
	Username        string      // Username for HTTP Basic Authentication.
	Password        string      // Password for HTTP Basic Authentication.
	APIKey          string      // Base64-encoded token for authorization; if set, overrides username/password and service token.
	Header          http.Header // Global HTTP request header.
	XSRFHeaderValue string      // Value for the kbn-xsrf header; defaults to "true" if not set.

	// PEM-encoded certificate authorities.
	// When set, an empty certificate pool will be created, and the certificates will be appended to it.
	// The option is only valid when the transport is not specified, or when it's http.Transport.
	CACert []byte

	RetryOnStatus []int                           // List of status codes for retry. Default: 502, 503, 504.
	DisableRetry  bool                            // Default: false.
	MaxRetries    int                             // Default: 3.
	RetryOnError  func(*http.Request, error) bool // Optional function allowing to indicate which error should be retried. Default: nil.

	EnableMetrics     bool // Enable the metrics collection.
	EnableDebugLogger bool // Enable the debug logging.

	RetryBackoff func(attempt int) time.Duration // Optional backoff duration. Default: nil.

	// Logger for client operations
	Transport http.RoundTripper         // The HTTP transport object.
	Logger    elastictransport.Logger   // The logger object.
	Selector  elastictransport.Selector // The selector object.

	Instrumentation elastictransport.Instrumentation // Enable instrumentation throughout the client.
}

type Client struct {
	Transport elastictransport.Interface
	//
	// API clients
	*kbapi.API

	// Configuration
	xsrfHeaderValue string

	// Internal state
	productCheckMu      sync.RWMutex
	productCheckSuccess bool
}

// NewOpenTelemetryInstrumentation provides the OpenTelemetry integration for Kibana client
func NewOpenTelemetryInstrumentation(provider trace.TracerProvider, captureBody bool) elastictransport.Instrumentation {
	return elastictransport.NewOtelInstrumentation(provider, captureBody, "1.0.0") // Replace with your version
}

// NewClient creates a new Kibana client
func NewClient(cfg Config) (*Client, error) {
	tp, err := newTransport(cfg)
	if err != nil {
		return nil, err
	}

	// Set default XSRF header value if not provided
	xsrfValue := cfg.XSRFHeaderValue
	if xsrfValue == "" {
		xsrfValue = "true"
	}

	client := &Client{
		Transport:       tp,
		xsrfHeaderValue: xsrfValue,
	}

	// Initialize API
	client.API = kbapi.New(client)

	return client, nil
}

// newTransport creates a new elastictransport client from configuration.
func newTransport(cfg Config) (*elastictransport.Client, error) {
	var addrs []string

	// Use provided addresses or environment variable
	if len(cfg.Addresses) == 0 {
		addrs = addrsFromEnvironment("KIBANA_URL")
	} else {
		addrs = append(addrs, cfg.Addresses...)
	}

	// Parse URLs
	urls, err := addrsToURLs(addrs)
	if err != nil {
		return nil, fmt.Errorf("cannot create client: %s", err)
	}

	// Default to localhost if no addresses provided
	if len(urls) == 0 {
		u, _ := url.Parse(defaultURL)
		urls = append(urls, u)
	}

	userAgent := initUserAgent()

	// Configure transport
	tpConfig := elastictransport.Config{
		UserAgent:         userAgent,
		URLs:              urls,
		Username:          cfg.Username,
		Password:          cfg.Password,
		APIKey:            cfg.APIKey,
		Header:            cfg.Header,
		CACert:            cfg.CACert,
		RetryOnStatus:     cfg.RetryOnStatus,
		DisableRetry:      cfg.DisableRetry,
		RetryOnError:      cfg.RetryOnError,
		MaxRetries:        cfg.MaxRetries,
		EnableMetrics:     cfg.EnableMetrics,
		EnableDebugLogger: cfg.EnableDebugLogger,
		Transport:         cfg.Transport,
		Logger:            cfg.Logger,
		Selector:          cfg.Selector,
		Instrumentation:   cfg.Instrumentation,
	}

	tp, err := elastictransport.New(tpConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating transport: %s", err)
	}

	return tp, nil
}

// InstrumentationEnabled propagates back to the client the Instrumentation provided by the transport.
func (c *Client) InstrumentationEnabled() elastictransport.Instrumentation {
	if tp, ok := c.Transport.(elastictransport.Instrumented); ok {
		return tp.InstrumentationEnabled()
	}
	return nil
}

// Perform delegates to Transport to execute a request and return a response.
func (c *Client) Perform(req *http.Request) (*http.Response, error) {
	req.Header.Set("kbn-xsrf", c.xsrfHeaderValue)

	// Perform the request
	res, err := c.Transport.Perform(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Metrics returns the client metrics.
func (c *Client) Metrics() (elastictransport.Metrics, error) {
	if mt, ok := c.Transport.(elastictransport.Measurable); ok {
		return mt.Metrics()
	}
	return elastictransport.Metrics{}, errors.New("transport is missing method Metrics()")
}

// addrsFromEnvironment returns a list of addresses by splitting
// the environment variable with comma, or an empty list.
func addrsFromEnvironment(envVar string) []string {
	var addrs []string
	if envURLs, ok := os.LookupEnv(envVar); ok && envURLs != "" {
		list := strings.Split(envURLs, ",")
		for _, u := range list {
			addrs = append(addrs, strings.TrimSpace(u))
		}
	}
	return addrs
}

// addrsToURLs creates a list of url.URL structures from url list.
func addrsToURLs(addrs []string) ([]*url.URL, error) {
	var urls []*url.URL
	for _, addr := range addrs {
		u, err := url.Parse(strings.TrimRight(addr, "/"))
		if err != nil {
			return nil, fmt.Errorf("cannot parse url: %v", err)
		}

		urls = append(urls, u)
	}
	return urls, nil
}

// initUserAgent creates the user agent string
func initUserAgent() string {
	var b strings.Builder
	b.WriteString("go-kibana")
	b.WriteRune('/')
	b.WriteString(Version)
	b.WriteRune(' ')
	b.WriteRune('(')
	b.WriteString(runtime.GOOS)
	b.WriteRune(' ')
	b.WriteString(runtime.GOARCH)
	b.WriteString("; ")
	b.WriteString("Go ")
	b.WriteString(runtime.Version())
	b.WriteRune(')')
	return b.String()
}

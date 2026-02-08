package extrator

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestHTTPServer_MetricAttributes(t *testing.T) {
	defaultRequest, err := http.NewRequest("GET", "http://example.com/path?query=test", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	var testCases = []struct {
		name             string
		givenRequest     *http.Request
		expectAttributes []attribute.KeyValue
	}{
		// FIXME add proper testcases
		{
			name:         "GET request, without body",
			givenRequest: defaultRequest,
			expectAttributes: []attribute.KeyValue{
				attribute.String("http.request.method", "GET"),
				attribute.String("url.scheme", "http"),
				attribute.String("server.address", "example.com"),
				attribute.String("network.protocol.name", "http"),
				attribute.String("network.protocol.version", "1.1"),
				attribute.Int64("http.response.status_code", 200),
				attribute.String("test", "test"),
			},
		},
		{
			name:         "server address",
			givenRequest: defaultRequest,
			expectAttributes: []attribute.KeyValue{
				attribute.String("http.request.method", "GET"),
				attribute.String("url.scheme", "http"),
				attribute.String("server.address", "example.com"),
				attribute.Int("server.port", 9999),
				attribute.String("network.protocol.name", "http"),
				attribute.String("network.protocol.version", "1.1"),
				attribute.Int64("http.response.status_code", 200),
				attribute.String("http.route", "/path/${id}"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ev := Values{}

			_ = ev.ExtractRequest(tc.givenRequest)

			attr := ev.MetricAttributes()

			assert.Len(t, attr, len(tc.expectAttributes))
			assert.ElementsMatch(t, tc.expectAttributes, attr)
		})
	}
}

func TestSplitAddress(t *testing.T) {
	var testCases = []struct {
		name          string
		whenHostPort  string
		expectHost    string
		expectPort    int
		expectedError string
	}{
		// FIXME add proper testcases

		//{"", "", -1},
		//{":8080", "", 8080},
		//{"127.0.0.1", "127.0.0.1", -1},
		//{"www.example.com", "www.example.com", -1},
		//{"127.0.0.1%25en0", "127.0.0.1%25en0", -1},
		//{"[]", "", -1}, // Ensure this doesn't panic.
		//{"[fe80::1", "", -1},
		//{"[fe80::1]", "fe80::1", -1},
		//{"[fe80::1%25en0]", "fe80::1%25en0", -1},
		//{"[fe80::1]:8080", "fe80::1", 8080},
		//{"[fe80::1]::", "", -1}, // Too many colons.
		//{"127.0.0.1:", "127.0.0.1", -1},
		//{"127.0.0.1:port", "127.0.0.1", -1},
		//{"127.0.0.1:8080", "127.0.0.1", 8080},
		//{"www.example.com:8080", "www.example.com", 8080},
		//{"127.0.0.1%25en0:8080", "127.0.0.1%25en0", 8080},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			host, port, err := SplitAddress(tc.whenHostPort)

			assert.Equal(t, tc.expectHost, host)
			assert.Equal(t, tc.expectPort, port)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHTTPMethod(t *testing.T) {
	var testCases = []struct {
		name                 string
		whenMethod           string
		expectMethod         string
		expectOriginalMethod string
	}{
		// FIXME add proper testcases

		//{"GET", "GET"},
		//{"get", "GET"},
		//{"POST", "POST"},
		//{"post", "POST"},
		//{"PUT", "PUT"},
		//{"put", "PUT"},
		//{"DELETE", "DELETE"},
		//{"delete", "DELETE"},
		//{"HEAD", "HEAD"},
		//{"head", "HEAD"},
		//{"OPTIONS", "OPTIONS"},
		//{"options", "OPTIONS"},
		//{"CONNECT", "CONNECT"},
		//{"connect", "CONNECT"},
		//{"TRACE", "TRACE"},
		//{"trace", "TRACE"},
		//{"PATCH", "PATCH"},
		//{"patch", "PATCH"},
		//{"unknown", "_OTHER"},
		//{"", "_OTHER"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			method, originalMethod := httpMethod(tc.whenMethod)

			assert.Equal(t, tc.expectMethod, method)
			assert.Equal(t, tc.expectOriginalMethod, originalMethod)
		})
	}
}

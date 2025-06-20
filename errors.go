package gofabric

import (
	"fmt"
)

// HTTPError represents an error returned by the Fabric API
type HTTPError struct {
	URL        string
	StatusCode int
	Body       *string
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	if e.Body == nil {
		return fmt.Sprintf("HTTP request to %s failed with status code: %d", e.URL, e.StatusCode)
	}
	return fmt.Sprintf(
		"HTTP request to %s failed with status code: %d: body: %s",
		e.URL,
		e.StatusCode,
		*e.Body,
	)
}

package types

import (
	"fmt"
	"strings"
)

// Headers //
type Headers map[string]string

func (h *Headers) String() string {
	return fmt.Sprint(*h)
}

// Set //
func (h *Headers) Set(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid header format, expected 'key: value': %s", value)
	}
	headerKey, headerValue := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	(*h)[headerKey] = headerValue
	return nil
}

// TestSpec //
type TestSpec struct {
	TotalRequests     int64
	ConcurrentThreads int64
	RequestTimeout    int64
	Method            string
	URL               string
	RequestHeaders    Headers
	Data              string
}

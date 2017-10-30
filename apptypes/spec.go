package apptypes

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/sirupsen/logrus"
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
	TotalRequests     uint64
	ConcurrentThreads uint64
	RequestTimeout    uint64
	Method            string
	URL               string
	RequestHeaders    Headers
	Data              string
}

// Print //
func (ts TestSpec) Print(logger *logrus.Logger) error {
	b, err := yaml.Marshal(ts)
	if err != nil {
		return fmt.Errorf("error printing spec: %s", err.Error())
	}
	logger.Info("spec:")
	for _, line := range strings.Split(string(b), "\n") {
		logger.Info("\t" + line)
	}
	return nil
}

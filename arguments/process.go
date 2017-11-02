package arguments

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

// ValidHTTPMethods //
var ValidHTTPMethods = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"DELETE",
	"CONNECT",
	"OPTIONS",
	"PATCH",
}

func validMethod(method string) bool {
	for _, m := range ValidHTTPMethods {
		if method == m {
			return true
		}
	}
	return false
}

// Process //
func Process(inputArgs []string, bodySource io.Reader, logger *logrus.Logger) (apptypes.TestSpec, logrus.Level, error) {
	var logLevelString string
	spec := apptypes.TestSpec{
		RequestHeaders: make(apptypes.Headers),
	}
	flagSet := flag.NewFlagSet("pressure", flag.ExitOnError)
	flagSet.StringVar(&logLevelString, "l", "info", "logging level")
	flagSet.Uint64Var(&spec.TotalRequests, "n", 100, "total number of requests (mininum 10)")
	flagSet.Uint64Var(&spec.ConcurrentThreads, "c", 10, "concurrent requests")
	flagSet.Uint64Var(&spec.RequestTimeout, "t", 30, "requests timeout")
	flagSet.StringVar(&spec.Data, "d", "", "data to be sent as body in request")
	flagSet.StringVar(&spec.Method, "X", "GET", "requests' HTTP method to use")
	flagSet.Var(&spec.RequestHeaders, "H", "header, can be specified multiple times")
	err := flagSet.Parse(inputArgs)
	if err != nil {
		return spec, logrus.InfoLevel, err
	}
	logLevel, err := logrus.ParseLevel(logLevelString)
	if err != nil {
		return spec,
			logrus.InfoLevel,
			fmt.Errorf(
				"invalid log level '%s', expecting one of %s",
				logLevelString,
				[]string{
					"panic",
					"fatal",
					"error",
					"warn",
					"info",
					"debug",
				},
			)
	}
	if flagSet.NArg() != 1 {
		return spec, logLevel, fmt.Errorf("one and only one URL argument must be specified")
	}
	if spec.TotalRequests <= 2 {
		return spec, logLevel, fmt.Errorf("total number of requets must be at least 3")
	}
	spec.URL = strings.TrimSpace(flagSet.Arg(0))
	spec.Method = strings.TrimSpace(strings.ToUpper(spec.Method))
	if !validMethod(spec.Method) {
		return spec, logLevel, fmt.Errorf(
			"invalid method '%s', valid methods are %s",
			spec.Method,
			strings.Join(ValidHTTPMethods, ", "),
		)
	}
	return spec, logLevel, nil
}

package attacker

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

func createHTTPClient(threadID uint64, spec apptypes.TestSpec, logger *logrus.Logger) *http.Client {
	logger.
		WithField("thread_id", threadID).
		Info("Setting up HTTP transport")
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	logger.
		WithField("thread_id", threadID).
		WithField("reuse_connections", spec.ReuseConnections).
		Info("Setting reuse of TCP/HTTP connections in transport")
	transport.DisableKeepAlives = !spec.ReuseConnections
	if spec.Proxy != "" {
		logger.
			WithField("thread_id", threadID).
			WithField("proxy", spec.Proxy).
			Info("Setting proxy")
		transport.Proxy = func(*http.Request) (*url.URL, error) {
			return url.Parse(spec.Proxy)
		}
	} else {
		logger.
			WithField("thread_id", threadID).
			Info("No proxy specified, using none")
	}
	logger.
		WithField("thread_id", threadID).
		Info("Creating client")
	client := &http.Client{}
	logger.
		WithField("thread_id", threadID).
		WithField("timeout_seconds", spec.RequestTimeout).
		Info("Setting timeout")
	client.Timeout = time.Duration(spec.RequestTimeout) * time.Second
	logger.
		WithField("thread_id", threadID).
		Info("Assigning HTTP transport to HTTP client")
	client.Transport = transport
	return client
}

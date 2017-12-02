package checkers

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

type accessorFunc func() error

type NetworkChecker struct {
	commonChecker

	URL *url.URL
}

type HTTPNetworkChecker struct {
	NetworkChecker
}

type TCPNetworkChecker struct {
	NetworkChecker
}

func (nc *NetworkChecker) access(results chan<- *CheckResult, accessor accessorFunc) {
	accessingURL := nc.URL.String()

	log.WithFields(log.Fields{
		"url":     accessingURL,
		"timeout": nc.Timeout,
	}).Debug("Access url")

	result := CheckResult{Ok: true, Producer: "url: " + accessingURL}
	if err := accessor(); err != nil {
		result.Ok = false
		result.Error = errors.Annotate(err, "Cannot access url "+accessingURL)
	}
	results <- &result
}

func (nc *NetworkChecker) isUnixPath() bool {
	return strings.HasPrefix(nc.URL.Host, "/")
}

func (hnc *HTTPNetworkChecker) Run(ctx context.Context, results chan<- *CheckResult) {
	hnc.access(results, func() error {
		req, err := http.NewRequest("GET", hnc.URL.String(), nil)
		if err != nil {
			return errors.Annotate(err, "Cannot compose request to the url")
		}

		newCtx, cancel := context.WithTimeout(ctx, hnc.Timeout)
		defer cancel()

		client := &http.Client{}
		if hnc.isUnixPath() {
			client = &http.Client{
				Transport: &http.Transport{
					DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
						return net.Dial("unix", hnc.URL.Host)
					},
				},
			}
		}

		response, err := client.Do(req.WithContext(newCtx))
		if err != nil {
			return errors.Annotate(err, "Cannot fetch URL data")
		}
		defer response.Body.Close()

		io.Copy(ioutil.Discard, response.Body)

		return nil
	})
}

func (tnc *TCPNetworkChecker) Run(ctx context.Context, results chan<- *CheckResult) {
	tnc.access(results, func() error {
		dialer := &net.Dialer{}

		network := "unix"
		if !tnc.isUnixPath() {
			network = tnc.URL.Scheme
		}

		newCtx, cancel := context.WithTimeout(ctx, tnc.Timeout)
		defer cancel()

		connection, err := dialer.DialContext(newCtx, network, tnc.URL.Host)
		if err != nil {
			return errors.Annotate(err, "Cannot dial")
		}
		defer connection.Close()

		return nil
	})
}

func NewNetworkChecker(timeout time.Duration, urlstruct *url.URL) (Checker, error) {
	switch urlstruct.Scheme {
	case "http", "https":
		return &HTTPNetworkChecker{
			NetworkChecker: NetworkChecker{commonChecker{timeout}, urlstruct},
		}, nil
	case "tcp", "udp":
		return &TCPNetworkChecker{
			NetworkChecker: NetworkChecker{commonChecker{timeout}, urlstruct},
		}, nil
	default:
		return nil, errors.Errorf("Unknown checker scheme %s", urlstruct.Scheme)
	}
}

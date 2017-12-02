package checkers

import (
	"net/url"
	"time"
)

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

type UDPNetworkChecker struct {
	NetworkChecker
}

func NewNetworkChecker(timeout time.Duration, urlstruct *url.URL) (Checker, error) {
	switch urlstruct.Scheme {
	case "http", "https":
		return &HTTPNetworkChecker{
			NetworkChecker: NetworkChecker{commonChecker{timeout}, urlstruct},
		}, nil
	case "tcp":
		return &TCPNetworkChecker{
			NetworkChecker: NetworkChecker{commonChecker{timeout}, urlstruct},
		}, nil
	case "udp":
		return &UDPNetworkChecker{
			NetworkChecker: NetworkChecker{commonChecker{timeout}, urlstruct},
		}, nil
	}
}

package config

import (
	"time"

	"github.com/juju/errors"
)

func validate(conf *Config) error {
	if conf.Periodicity.Duration <= time.Duration(0) {
		return errors.New("Periodicity should be >= 0")
	}
	if conf.Host == "" {
		return errors.New("HTTP endpoint host should be set")
	}
	if conf.Port == 0 {
		return errors.New("HTTP endpoint port should be set")
	}
	if conf.OkStatus < 100 || conf.OkStatus >= 600 {
		return errors.New("HTTP status code is incorrect")
	}
	if conf.NokStatus < 100 || conf.NokStatus >= 600 {
		return errors.New("HTTP status code is incorrect")
	}

	for _, check := range conf.Checks {
		if check.Timeout.Duration <= time.Duration(0) {
			return errors.New("Checker timeout >= 0")
		}

		var err error
		switch check.Type {
		case "network":
			err = validateNetwork(&check)
		case "command", "shell":
			err = validateCommand(&check)
		default:
			return errors.Errorf("Unknown checker type %s", check.Type)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func validateNetwork(check *ConfigChecker) error {
	if check.Exec != "" {
		return errors.New("Cannot set exec for network checker")
	}
	if check.URL.URL == nil {
		return errors.New("URL has to be defined for network checker")
	}

	for _, v := range check.StatusCodes {
		if v < 100 || v >= 600 {
			return errors.New("Unknown HTTP status code")
		}
	}

	switch check.URL.Scheme {
	case "http", "https", "tcp", "udp":
		return nil
	default:
		return errors.Errorf("Unknown URL scheme %s", check.URL.Scheme)
	}
}

func validateCommand(check *ConfigChecker) error {
	if check.Exec == "" {
		return errors.New("Command checker must have exec attribute")
	}
	if check.URL.URL != nil {
		return errors.New("Command checker must not have url attribute")
	}

	return nil
}

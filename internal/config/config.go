package config

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
)

type duration struct {
	time.Duration
}

func (dur *duration) UnmarshalText(text []byte) (err error) {
	dur.Duration, err = time.ParseDuration(string(text))
	return
}

type urltype struct {
	*url.URL
}

func (ut *urltype) UnmarshalText(text []byte) (err error) {
	chunks := strings.SplitN(string(text), "//", 2)

	if len(chunks) != 2 {
		return errors.New("Unknown scheme")
	}
	scheme, newUrl := chunks[0], chunks[1]

	chunks = strings.SplitN(newUrl, "/", 2)
	hostname := chunks[0]
	rest := ""
	if len(chunks) == 2 {
		rest = chunks[1]
		value, err := url.PathUnescape(hostname)
		if err != nil {
			return err
		}
		hostname = value
		rest = "/" + rest
	}

	newUrl = fmt.Sprintf("%s//localhost:80%s", scheme, rest)
	ut.URL, err = url.Parse(newUrl)
	if err != nil {
		return errors.Trace(err)
	}

	ut.URL.Host = hostname

	return nil
}

type Config struct {
	Host        string
	Port        uint16
	Periodicity duration
	Checks      []ConfigChecker
}

type ConfigChecker struct {
	Type    string
	Exec    string
	URL     urltype
	Timeout duration
}

func ParseConfigFile(file *os.File) (*Config, error) {
	conf := &Config{}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Annotate(err, "Cannot read from the config file")
	}

	_, err = toml.Decode(string(buf), conf)
	if err != nil {
		return nil, errors.Annotate(err, "Cannot parse config file")
	}

	for idx := range conf.Checks {
		if conf.Checks[idx].Timeout.Duration == time.Duration(0) {
			conf.Checks[idx].Timeout.Duration = conf.Periodicity.Duration
		}
	}

	err = validate(conf)
	if err != nil {
		return nil, errors.Annotate(err, "Cannot validate config file")
	}

	return conf, nil
}

package nca

import (
	"fmt"
	"github.com/kballard/go-shellquote"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

// Config describes the full agent configuration.
type Config struct {
	Publishers map[string]map[string]interface{}
	Hostname   string
	Checks     map[string]Check
}

// ReadConfig loads configuration from the given source and returns a
// fully initialized Configuration struct from it.
func ReadConfig(src io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	for name, check := range c.Checks {
		if check.Name == "" {
			check.Name = name
		}
		if check.Interval < 1 {
			check.Interval = 60
		}
		if check.Retry < 1 {
			check.Retry = 60
		}
		if check.Timeout < 1 {
			check.Timeout = 10
		}

		splitArgs, err := shellquote.Split(check.Command)
		if err != nil {
			return nil, err
		}
		if len(splitArgs) < 1 {
			return nil, Error{
				Code:    ErrCheckMissingCommand,
				Message: fmt.Sprintf("Check '%s' is missing a command to execute", name),
			}
		}
		check.Args = splitArgs

		c.Checks[name] = check
	}

	return c, nil
}

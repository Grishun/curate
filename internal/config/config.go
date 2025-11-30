package config

import (
	"time"

	cli "github.com/urfave/cli/v3"
)

type Config struct {
	HTTPPort       string        `yaml:"http_port"`
	WorkerInterval time.Duration `yaml:"worker_interval"`
	Currencies     []string      `yaml:"currencies"`
	Quote          string        `yaml:"quote"`
	CoindeskURL    string        `yaml:"coindesk_url"`
	CoindeskToken  string        `yaml:"coindesk_token"`
	//TODO: add limit
}

// New parse config from urfave/cli flags
func New(c *cli.Command) *Config {
	cfg := Config{
		HTTPPort:       c.String("http-port"),
		WorkerInterval: c.Duration("provider-interval"),
		Currencies:     c.StringSlice("currencies"),
		Quote:          c.String("quote"),
		CoindeskURL:    c.String("coindesk-url"),
		CoindeskToken:  c.String("coindesk-token"),
	}

	return &cfg
}

// NewFromYAML
// TODO: implement
func NewFromYAML() {}

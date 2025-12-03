package config

import (
	"time"

	cli "github.com/urfave/cli/v3"
)

type Config struct {
	RestHost        string        `yaml:"rest_host"`
	RestPort        string        `yaml:"rest_port"`
	PollingInterval time.Duration `yaml:"polling_interval"`
	Currencies      []string      `yaml:"currencies"`
	Quote           string        `yaml:"quote"`
	CoindeskURL     string        `yaml:"coindesk_url"`
	CoindeskToken   string        `yaml:"coindesk_token"`
	HistoryLimit    uint32        `yaml:"history_limit"`
	InfluxDBURI     string        `yaml:"influxdb_uri"`
	InfluxDBToken   string        `yaml:"influxdb_token"`
	InfluxDBBucket  string        `yaml:"influxdb_bucket"`
	InMemoryStorage bool          `yaml:"in_memory_storage"`
}

// New parse config from urfave/cli flags and create Config struct
func New(c *cli.Command) *Config {
	cfg := Config{
		RestHost:        c.String("rest-host"),
		RestPort:        c.String("rest-port"),
		PollingInterval: c.Duration("polling-interval"),
		Currencies:      c.StringSlice("currencies"),
		Quote:           c.String("quote"),
		CoindeskURL:     c.String("coindesk-url"),
		CoindeskToken:   c.String("coindesk-token"),
		HistoryLimit:    c.Uint32("history-limit"),
		InfluxDBURI:     c.String("influxdb-uri"),
		InfluxDBToken:   c.String("influxdb-token"),
		InfluxDBBucket:  c.String("influxdb-bucket"),
		InMemoryStorage: c.Bool("in-memory-storage"),
	}

	return &cfg
}

// NewFromYAML
// TODO: implement
func NewFromYAML() {}

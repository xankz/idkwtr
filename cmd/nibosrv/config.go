package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	// Host that the server will listen on.
	Host string `json:"host"`

	// Port that the server will listen to.
	Port int `json:"port"`

	// StaticFilesDir is the path to the static files directory.
	StaticFilesDir string `json:"static_dir"`
}

func loadConfig() (config, error) {
	// Read flag configuration first, then override with file configuration
	// (if config file specified).
	conf := parseFlagConfig()
	if *flagConfig != "" {
		fb, err := os.ReadFile(*flagConfig)
		if err != nil {
			return config{}, fmt.Errorf("reading config file: %w", err)
		}
		if err := json.Unmarshal(fb, &conf); err != nil {
			return config{}, fmt.Errorf("parsing config file: %w", err)
		}
	}

	return conf, nil
}

func parseFlagConfig() config {
	return config{
		Host: *flagHost,
		Port: *flagPort,
	}
}

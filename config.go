package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	log "log/slog"

	"gopkg.in/yaml.v3"
)

type logConfig struct {
	Level string `yaml:"level"`
}

type exporterConfig struct {
	Listen           string        `yaml:"listen"`
	Path             string        `yaml:"path"`
	Namespace        string        `yaml:"namespace"`
}

type UbntConfig struct {
	SkipSslVlidation bool   `yaml:"skip_ssl_validation"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
}

// Config -
type Config struct {
	Log      logConfig      `yaml:"log"`
	Exporter exporterConfig `yaml:"exporter"`
	Ubnt     UbntConfig     `yaml:"ubnt"`
}

func (c *exporterConfig) validate() error {
	if 0 == len(c.Path) {
		c.Path = "/metrics"
	}
	if 0 == len(c.Listen) {
		return fmt.Errorf("missing key 'expoerter.listen'")
	}
	return nil
}

func (c *UbntConfig) validate() error {
	if len(c.Username) == 0 {
		return fmt.Errorf("missing mandatory key ubnt.username")
	}
	if len(c.Password) == 0 {
		return fmt.Errorf("missing mandatory key ubnt.password")
	}
	return nil
}

// Validate - Validate configuration object
func (c *Config) Validate() error {
	if err := c.Ubnt.validate(); err != nil {
		return fmt.Errorf("invalid ubnt configuration: %s", err)
	}
	if err := c.Exporter.validate(); err != nil {
		return fmt.Errorf("invalid exporter configuration: %s", err)
	}
	return nil
}

// NewConfig - Creates and validates config from given reader
func NewConfig(file io.Reader) *Config {
	content, err := io.ReadAll(file)
	if err != nil {
		log.Error("unable to read configuration file", "error", err)
		os.Exit(1)
	}
	config := Config{}
	if err = yaml.Unmarshal(content, &config); err != nil {
		log.Error("unable to read configuration yaml file", "error", err)
		if err = json.Unmarshal(content, &config); err != nil {
			log.Error("unable to read configuration json file", "error", err)
			os.Exit(1)
		}
	}
	if err = config.Validate(); err != nil {
		log.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	return &config
}

package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Zabbix   ZabbixConfig `toml:"zabbix"`
	OTLP     OTLPConfig   `toml:"otlp"`
	Salt     SaltConfig   `toml:"salt"`
	LogLevel string       `toml:"log_level"`
}

type ZabbixConfig struct {
	URL      string `toml:"url"`
	Token    string `toml:"token"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Timeout  int    `toml:"timeout"`
}

type SaltConfig struct {
	URL      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	EAuth    string `toml:"eauth"`
}

type OTLPConfig struct {
	Endpoint    string `toml:"endpoint"`
	Protocol    string `toml:"protocol"`
	ServiceName string `toml:"service_name"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Default values
	if cfg.Zabbix.Timeout == 0 {
		cfg.Zabbix.Timeout = 30
	}
	if cfg.OTLP.ServiceName == "" {
		cfg.OTLP.ServiceName = "zabbix-dna"
	}

	return &cfg, nil
}



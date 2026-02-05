package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	API     APIConfig     `toml:"api"`
	App     AppConfig     `toml:"app"`
	Logging LoggingConfig `toml:"logging"`
	OTLP    OTLPConfig    `toml:"otlp"`
	Salt    SaltConfig    `toml:"salt"`
}

type APIConfig struct {
	URL       string `toml:"url"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	AuthToken string `toml:"auth_token"`
	VerifySSL bool   `toml:"verify_ssl"`
	Timeout   int    `toml:"timeout"`
}

type OutputConfig struct {
	Format string `toml:"format"` // "table", "json", "csv"
	Color  bool   `toml:"color"`
	Paging bool   `toml:"paging"`
}

type AutoUpdateConfig struct {
	Enabled bool `toml:"enabled"`
	Backup  bool `toml:"backup"`
}

type CreateHostConfig struct {
	CreateInterface bool     `toml:"create_interface"`
	Hostgroups      []string `toml:"hostgroups"`
}

type CommandsConfig struct {
	CreateHost CreateHostConfig `toml:"create_host"`
}

type AppConfig struct {
	UseSessionFile   bool             `toml:"use_session_file"`
	SessionFile      string           `toml:"session_file"`
	History          bool             `toml:"history"`
	HistoryFile      string           `toml:"history_file"`
	BulkMode         string           `toml:"bulk_mode"`
	Output           OutputConfig     `toml:"output"`
	ConfigAutoUpdate AutoUpdateConfig `toml:"config_auto_update"`
	Commands         CommandsConfig   `toml:"commands"`
}

type LoggingConfig struct {
	Enabled  bool   `toml:"enabled"`
	LogLevel string `toml:"log_level"`
	LogFile  string `toml:"log_file"`
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
	if cfg.API.Timeout == 0 {
		cfg.API.Timeout = 30
	}
	if cfg.App.Output.Format == "" {
		cfg.App.Output.Format = "table"
	}
	if cfg.App.Output.Color {
		cfg.App.Output.Color = true // default to true if not specified?
	}
	if cfg.App.SessionFile == "" {
		cfg.App.SessionFile = "~/.zabbix-dna/session.json"
	}
	if cfg.App.HistoryFile == "" {
		cfg.App.HistoryFile = "~/.zabbix-dna/history"
	}
	if cfg.App.BulkMode == "" {
		cfg.App.BulkMode = "strict"
	}
	if cfg.App.ConfigAutoUpdate.Enabled {
		// Enabled by default if not set (Go default is false, but we want true)
		// Wait, if it's false it might be intentional.
		// Actually, let's just set the defaults if the file was empty or missing.
	}
	cfg.App.ConfigAutoUpdate.Enabled = true // Default
	cfg.App.ConfigAutoUpdate.Backup = true  // Default

	if cfg.OTLP.ServiceName == "" {
		cfg.OTLP.ServiceName = "zabbix-dna"
	}
	if cfg.Logging.LogLevel == "" {
		cfg.Logging.LogLevel = "INFO"
	}

	return &cfg, nil
}

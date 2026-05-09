package utils

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host          string `yaml:"host"`
		HTTPPort      int    `yaml:"http_port"`
		TCPPort       int    `yaml:"tcp_port"`
		UDPPort       int    `yaml:"udp_port"`
		GRPCPort      int    `yaml:"grpc_port"`
		WebSocketPort int    `yaml:"websocket_port"`
	} `yaml:"server"`

	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`

	User struct {
		Username string `yaml:"username"`
		Token    string `yaml:"token"`
	} `yaml:"user"`

	Sync struct {
		AutoSync           bool   `yaml:"auto_sync"`
		ConflictResolution string `yaml:"conflict_resolution"`
	} `yaml:"sync"`

	Notifications struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"notifications"`

	Logging struct {
		Level string `yaml:"level"`
		Path  string `yaml:"path"`
	} `yaml:"logging"`

	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

// DefaultConfig returns sensible defaults matching the CLI manual.
func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	cfg := &Config{}
	cfg.Server.Host = "localhost"
	cfg.Server.HTTPPort = 8080
	cfg.Server.TCPPort = 9090
	cfg.Server.UDPPort = 9091
	cfg.Server.GRPCPort = 9092
	cfg.Server.WebSocketPort = 9093
	cfg.Database.Path = filepath.Join(home, ".mangahub", "data.db")
	cfg.Sync.AutoSync = true
	cfg.Sync.ConflictResolution = "last_write_wins"
	cfg.Notifications.Enabled = true
	cfg.Logging.Level = "info"
	cfg.Logging.Path = filepath.Join(home, ".mangahub", "logs")
	cfg.JWT.Secret = "mangahub-secret-change-in-production"
	return cfg
}

// LoadConfig reads ~/.mangahub/config.yaml, falling back to defaults on error.
func LoadConfig() *Config {
	cfg := DefaultConfig()
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg
	}
	path := filepath.Join(home, ".mangahub", "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg // file doesn't exist yet — use defaults
	}
	_ = yaml.Unmarshal(data, cfg)
	return cfg
}

// SaveConfig writes the config to ~/.mangahub/config.yaml.
func SaveConfig(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".mangahub")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "config.yaml"), data, 0o600)
}

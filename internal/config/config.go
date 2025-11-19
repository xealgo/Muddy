package config

import (
	"crypto/tls"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Application configuration
type Config struct {
	Port      int
	CertFile  string
	KeyFile   string
	TLSConfig *tls.Config
	envPath   string
}

type ConfigOption func(*Config)

// Creates a new Config instance
func NewConfig(opts ...ConfigOption) (*Config, error) {
	cfg := &Config{
		Port:     8080,
		CertFile: "server.crt",
		KeyFile:  "server.key",
		envPath:  ".env",
	}

	for _, opts := range opts {
		opts(cfg)
	}

	err := cfg.LoadFromEnv()
	if err != nil {
		return nil, err
	}

	err = cfg.SetFromEnv()
	if err != nil {
		return nil, err
	}

	tlsConfig, err := cfg.LoadCerts()
	if err != nil {
		return nil, err
	}

	cfg.TLSConfig = tlsConfig

	return cfg, nil
}

// WithEnvPath sets the custom .env path
func WithEnvPath(path string) ConfigOption {
	return func(cfg *Config) {
		cfg.envPath = path
	}
}

// WithDefaults sets default configuration values
func WithDefaults(port int, certFile string, keyFile string) ConfigOption {
	return func(cfg *Config) {
		cfg.Port = port
		cfg.CertFile = certFile
		cfg.KeyFile = keyFile
	}
}

// Loads configuration from a .env file
func (cfg *Config) LoadFromEnv() error {
	// Check if file exists
	_, err := os.Stat(cfg.envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ConfigError{Type: FileNotFound, Message: "Config file not found", EnvPath: cfg.envPath, Wrapped: err}
		}
	}

	err = godotenv.Load(cfg.envPath)
	if err != nil {
		return ConfigError{Message: "Failed to load environment file", Wrapped: err, EnvPath: cfg.envPath}
	}

	return nil
}

// Loads TLS certificates
func (cfg *Config) LoadCerts() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, ConfigError{Type: LoadFailure, Message: "Failed to load certificate", CertFile: cfg.CertFile, KeyFile: cfg.KeyFile, Wrapped: err}
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"wq-vvv-01", "h3", "h2", "http/1.1"},
	}, nil
}

// Sets configuration values from environment variables
func (cfg *Config) SetFromEnv() error {
	envPort := GetEnv("PORT", "")

	if envPort != "" {
		port, err := strconv.Atoi(envPort)
		if err != nil {
			return ConfigError{Type: InvalidValue, Message: "Invalid port value", EnvPath: cfg.envPath, Wrapped: err}
		}

		cfg.Port = port
	}

	cfg.CertFile = GetEnv("CERT_FILE", cfg.CertFile)

	_, err := os.Stat(cfg.CertFile)
	if err != nil {
		return ConfigError{Type: FileNotFound, Message: "Certificate file not found", EnvPath: cfg.envPath, CertFile: cfg.CertFile, Wrapped: err}
	}

	cfg.KeyFile = GetEnv("KEY_FILE", cfg.KeyFile)

	_, err = os.Stat(cfg.KeyFile)
	if err != nil {
		return ConfigError{Type: FileNotFound, Message: "Key file not found", EnvPath: cfg.envPath, KeyFile: cfg.KeyFile, Wrapped: err}
	}

	return nil
}

// Attempts to load an env by name or returns the default value if not provided
func GetEnv(name string, defValue string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	return defValue
}

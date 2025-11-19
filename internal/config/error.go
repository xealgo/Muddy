package config

import (
	"fmt"
	"strings"
)

type ConfigErrorType string

const (
	LoadFailure  ConfigErrorType = "LOAD_FAILURE"
	FileNotFound ConfigErrorType = "FILE_NOT_FOUND"
	InvalidValue ConfigErrorType = "INVALID_VALUE"
)

// ConfigError represents an error in configuration loading
type ConfigError struct {
	Type     ConfigErrorType
	Message  string
	EnvPath  string
	CertFile string
	KeyFile  string
	Wrapped  error
}

// NewConfigError creates a new ConfigError instance
func (e ConfigError) Error() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Config error: %v, %s - ", e.Type, e.Message))

	if e.EnvPath != "" {
		builder.WriteString(fmt.Sprintf(" (env: %s)", e.EnvPath))
	}

	if e.CertFile != "" {
		builder.WriteString(fmt.Sprintf(" (cert: %s)", e.CertFile))
	}

	if e.KeyFile != "" {
		builder.WriteString(fmt.Sprintf(" (key: %s)", e.KeyFile))
	}

	builder.WriteString(fmt.Sprintf(" - %v", e.Wrapped))

	return builder.String()
}

// Unwrap returns the underlying error
func (e ConfigError) Unwrap() error {
	return e.Wrapped
}

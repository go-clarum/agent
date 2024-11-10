package config

import (
	"log/slog"
	"strings"
)

const (
	version              string = "1.0.0-snapshot"
	defaultBaseDir       string = "."
	defaultConfigFile    string = "clarum-properties.yaml"
	defaultLogLevel      string = "info"
	defaultProfile       string = "dev"
	defaultActionTimeout uint   = 10
	defaultAgentPort     uint   = 9091
)

// replace missing attributes from the configuration with their default values
func (config *config) setDefaults() {
	slog.Debug("Replacing missing values with defaults")

	if strings.TrimSpace(config.Profile) == "" {
		config.Profile = defaultProfile
	}
	if strings.TrimSpace(config.Logging.Level) == "" {
		config.Logging.Level = defaultLogLevel
	}
	if config.Actions.TimeoutSeconds == 0 {
		config.Actions.TimeoutSeconds = defaultActionTimeout
	}
	if config.Agent.Port == 0 {
		config.Agent.Port = defaultAgentPort
	}
}

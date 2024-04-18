package config

import (
	"github.com/go-clarum/agent/files"
	"github.com/go-clarum/agent/logging"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path"
	"time"
)

var c *config

type config struct {
	Agent struct {
		Port uint
	}
	Profile string
	Actions struct {
		TimeoutSeconds uint
	}
	Logging struct {
		Level string
	}
}

func LoadConfig() {
	configFilePath := path.Join(*baseDir, *configFile)
	conf, err := files.ReadYamlFileToStruct[config](configFilePath)
	if err != nil {
		dir, _ := os.Getwd()
		logging.Infof("No config file found in [%s] - default values will be used instead", dir)
		conf = &config{}
	}

	conf.setDefaults()
	conf.overwriteWithCliFlags()
	c = conf

	configYaml, _ := yaml.Marshal(conf)
	logging.Infof("Using the following config:\n[\n%s]", configYaml)
}

func Version() string {
	return version
}

func BaseDir() string {
	return *baseDir
}

func LoggingLevel() slog.Level {
	return logging.ParseLevel(c.Logging.Level)
}

func ActionTimeout() time.Duration {
	return time.Duration(c.Actions.TimeoutSeconds) * time.Second
}

func AgentPort() uint {
	return c.Agent.Port
}

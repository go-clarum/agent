package config

import (
	"github.com/go-clarum/agent/infrastructure/files"
	"gopkg.in/yaml.v3"
	"log"
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

func init() {
	initLogger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

	configFilePath := path.Join(*baseDir, *configFile)
	conf, err := files.ReadYamlFileToStruct[config](configFilePath)
	if err != nil {
		dir, _ := os.Getwd()
		initLogger.Printf("No config file found in [%s] - default values will be used instead", dir)
		conf = &config{}
	}

	conf.setDefaults()
	conf.overwriteWithCliFlags()
	c = conf

	configYaml, _ := yaml.Marshal(conf)
	initLogger.Printf("Using the following config:\n[\n%s]", configYaml)
}

func Version() string {
	return version
}

func BaseDir() string {
	return *baseDir
}

func LoggingLevel() string {
	return c.Logging.Level
}

func ActionTimeout() time.Duration {
	return time.Duration(c.Actions.TimeoutSeconds) * time.Second
}

func AgentPort() uint {
	return c.Agent.Port
}

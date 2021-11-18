package config

import (
	"github.com/Mortimor1/telegraf-discovery/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Debug *bool `yaml:"debug"`
	Jobs  []Job `yaml:"jobs"`
}

type Job struct {
	Subnet        string `yaml:"subnet"`
	ConfigFile    string `yaml:"config_file"`
	ContainerName string `yaml:"container_name"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Read Config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config/config.yml", instance); err != nil {
			desc, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(desc)
			logger.Fatal(err)
		}
	})
	return instance
}

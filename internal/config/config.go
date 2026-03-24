package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string        `yaml:"env" env-default:"8080"`
	ServerCfg  ServerConfig  `yaml:"server"`
	StorageCfg StorageConfig `yaml:"storage"`
}

type ServerConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"10s"`
}

type StorageConfig struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBname   string `yaml:"dbname" env-required:"true"`
	Port     int    `yaml:"port" env-default:"5432"`
}

func MustLoadConfig(path string) Config {
	if _, err := os.Stat(path); err != nil {
		panic(err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return cfg
}

package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env        string        `yaml:"env" env-default:"8080"`
	ServerCfg  ServerConfig  `yaml:"server"`
	StorageCfg StorageConfig `yaml:"postgres"`
	RedisCfg   RedisConfig   `yaml:"redis"`
}

type ServerConfig struct {
	Port    int           `yaml:"port" env:"PORT_SERVER" env-required:"true"`
	Host    string        `yaml:"host" env:"HOST_SERVER" env-default:"localhost"`
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT_SERVER" env-default:"10s"`
}

type StorageConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Username string `yaml:"username" env:"POSTGRES_USERNAME" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	DBname   string `yaml:"dbname" env:"POSTGRES_DBNAME" env-required:"true"`
	Port     int    `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
}

type RedisConfig struct {
	Host     string        `yaml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Port     int           `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
	User     string        `yaml:"user" env:"REDIS_USER" env-default:"admin"`
	Password string        `yaml:"password" env:"REDIS_PASSWORD" env-default:"admin"`
	DB       int           `yaml:"db" env:"REDIS_DB" env-default:"0"`
	TTL      time.Duration `yaml:"ttl" env:"REDIS_TTL" env-default:"60s"`
}

// MustLoadConfig загружает конфиг из .env файлов и системных переменных
// Приоритет: системный ENV > .env файл > значения по умолчанию
func MustLoadConfig(filenames ...string) Config {
	err := godotenv.Load(filenames...)
	if err != nil {
		panic(err.Error())
	}

	return Config{
		Env:        getEnvString("APP_ENV", "development"), // добавьте APP_ENV в .env
		ServerCfg:  loadServerConfig(),
		StorageCfg: loadStorageConfig(),
		RedisCfg:   loadRedisConfig(),
	}
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Host:    getEnvString("HOST_SERVER", "localhost"),
		Port:    getEnvInt("PORT_SERVER", 8080),
		Timeout: getEnvDuration("TIMEOUT_SERVER", 10*time.Second),
	}
}

func loadStorageConfig() StorageConfig {
	return StorageConfig{
		Host:     getEnvString("POSTGRES_HOST", "localhost"),
		Username: getEnvString("POSTGRES_USERNAME", ""),
		Password: getEnvString("POSTGRES_PASSWORD", ""),
		DBname:   getEnvString("POSTGRES_DBNAME", ""),
		Port:     getEnvInt("POSTGRES_PORT", 5432),
	}
}

func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     getEnvString("REDIS_HOST", "localhost"),
		Port:     getEnvInt("REDIS_PORT", 6379),
		User:     getEnvString("REDIS_USER", "admin"),
		Password: getEnvString("REDIS_PASSWORD", "admin"),
		DB:       getEnvInt("REDIS_DB", 0),
		TTL:      getEnvDuration("REDIS_TTL", 60*time.Second),
	}
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if durationValue, err := time.ParseDuration(value); err == nil {
			return durationValue
		}
	}
	return defaultValue
}

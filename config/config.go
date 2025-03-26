package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	Env string `yaml:"env" env:"ENV"`

	DB struct {
		DSN      string `yaml:"dsn" env:"DB_DSN"`
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		Database string `yaml:"database" env:"DB_NAME"`
	}

	GRPC struct {
		Port int `yaml:"port" env:"GRPC_PORT" envDefault:"50051"`
	}
}

func MustLoad(logger *zap.Logger) *Config {
	var cfg Config

	path := fetchConfig()

	logger.Info("using config file", zap.String("path", path))

	if path == "" {
		// Загрузка переменных из .env файла
		if err := godotenv.Load(".env"); err != nil {
			logger.Error("failed to load .env file", zap.Error(err))
		}

		if err := envconfig.Process("", &cfg); err != nil {
			logger.Error("failed to process envconfig", zap.Error(err))
			panic(err)
		}

		logger.Info("using config by env", zap.Any("config", cfg))
		return &cfg
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist")
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	logger.Info("using config by flag", zap.Any("config", cfg))

	return &cfg
}

func fetchConfig() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	return res
}

package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env"`
	StoragePath string        `yaml:"storage_path"`
	TokenTTL    time.Duration `yaml:"token_ttl"`
	GRPC        `yaml:"grpc"`
}

type GRPC struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func fetchConfigFlag() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "Path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res

}

func MustLoad() *Config {
	// os.Setenv("CONFIG_PATH", "./config/local.yaml")

	conf_path := fetchConfigFlag()
	if conf_path == "" {
		log.Fatal("CONFIG_PATH is empty")
	}

	if _, err := os.Stat(conf_path); os.IsNotExist(err) {
		log.Fatalf("Config file %s does not exist", conf_path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(conf_path, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}

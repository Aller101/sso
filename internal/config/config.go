package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yaml:"env"`
	Storage_path string `yaml:"storage_path"`
	GRPC         `yaml:"grpc"`
}

type GRPC struct {
	Port    int `yaml:"port"`
	Timeout int `yaml:"timeout"`
}

func LoadFlag() {

}

func MustLoad() *Config {
	// os.Setenv("CONFIG_PATH", "./config/local.yaml")

	conf_path := os.Getenv("CONFIG_PATH")
	if conf_path == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat("CONFIG_PATH"); os.IsNotExist(err) {
		log.Fatalf("Config file %s does not exist", conf_path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(conf_path, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}

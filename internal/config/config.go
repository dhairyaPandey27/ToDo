package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string
}

// This should be set as if our env-required was false then this default value can be taken as some other thing, then it would be a problem
//env-default: "production"

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePort string `yaml:"storage_port" env-required:"true"`
	HTTPServer  `yaml:"http-server"`
}

func MustLoad() *Config {
	// Get the configuration file path from the CONFIG_PATH environment variable
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	// If CONFIG_PATH is not set, try to get the path from the command-line flag
	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		// Dereference the flag pointer to get the actual path
		configPath = *flags

		// If neither the environment variable nor the flag is set, stop the program
		if configPath == "" {
			log.Fatal("Config Path is not set")
		}
	}

	// Check whether the configuration file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	// Create an empty Config struct
	var cfg Config

	// Read the configuration file and store its values in cfg
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		// Stop the program if the configuration file cannot be read
		log.Fatalf("can not read config file: %s", err.Error())
	}

	// Return a pointer to the loaded configuration
	return &cfg
}

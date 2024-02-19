package config

import (
	"errors"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Bot BotConf `yaml:"bot"`
	DB  DBConf  `yaml:"db-conn"`
}

type BotConf struct {
	Token  string  `yaml:"token"`
	Admins []int64 `yaml:"admins,flow"`
}

type DBConf struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db-name"`
}

func InitConfig() (*Config, error) {
	config := &Config{}

	cfgPath, err := getEnv("CONFIG_PATH")
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func getEnv(envKey string) (string, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Printf("err loading: %v\n", err)
	}

	env := os.Getenv(envKey)
	if env == "" {
		return "", errors.New("missing addres")
	}

	return env, nil
}

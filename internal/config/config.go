package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     Server     `yaml:"server"`
	PostgresDB PostgresDB `yaml:"postgres"`
	Logger     Logger     `yaml:"logger"`
	Redis      Redis      `yaml:"redis"`
	JWT        JWT        `yaml:"jwt"`
}

type Logger struct {
	Level  int    `yaml:"level"`
	Format string `yaml:"format"`
}

type PostgresDB struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	SslMode      string `yaml:"ssl_mode"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxLifeTime  int    `yaml:"max_life_time"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Server struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

func LoadConfig(path string) (*Config, error) {
	confFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err = yaml.Unmarshal(confFile, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

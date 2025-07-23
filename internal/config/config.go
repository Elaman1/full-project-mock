package config

import (
	"github.com/caarlos0/env/v10"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     Server     `yaml:"server"`
	PostgresDB PostgresDB `yaml:"postgres"`
	Logger     Logger     `yaml:"logger"`
	Redis      Redis      `yaml:"redis"`
	JWT        JWTConfig  `yaml:"jwt"`
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

type JWTConfig struct {
	PrivateKeyPath string `env:"JWT_PRIVATE_KEY_PATH,required"`
	PublicKeyPath  string `env:"JWT_PUBLIC_KEY_PATH,required"`
	AccessTTL      string `yaml:"access_ttl"`
}

func LoadConfig(path, envPath string) (*Config, error) {
	confFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(confFile, &cfg); err != nil {
		return nil, err
	}

	if err = godotenv.Load(envPath); err != nil {
		return nil, err
	}

	// Переопределим значениями из ENV, если они заданы
	if err = env.Parse(&cfg.JWT); err != nil {
		return nil, err
	}

	if err = validateCfg(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

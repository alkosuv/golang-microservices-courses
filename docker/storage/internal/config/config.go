package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"path/filepath"
)

type Config struct {
	App struct {
		HTTP Service `yaml:"http"`
	} `yaml:"app"`

	Database struct {
		Host     string `yaml:"host" env:"DATABASE_HOST"`
		Port     uint16 `yaml:"port" env:"DATABASE_PORT"`
		Database string `yaml:"database" env:"DATABASE_DATABASE"`
		User     string `yaml:"user" env:"DATABASE_USER"`
		Password string `yaml:"password" env:"DATABASE_PASSWORD"`
		SSLmode  string `yaml:"sslmode" env:"DATABASE_SSLMODE"`
	} `yaml:"database"`
}

type Service struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func NewConfig(configPath string) (config Config, err error) {
	filename, err := filepath.Abs(configPath)
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}

	if err := cleanenv.ParseYAML(file, &config); err != nil {
		return Config{}, err
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		return Config{}, err
	}

	return
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User, c.Database.Password,
		c.Database.Host, c.Database.Port,
		c.Database.Database,
		c.Database.SSLmode,
	)
}

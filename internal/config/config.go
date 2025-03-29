package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type (
	Config struct {
		Env     string  `yaml:"env" env:"ENV" env-default:"local"`
		App     App     `yaml:"app"`
		Log     Log     `yaml:"log"`
		PG      PG      `yaml:"pg"`
		HTTP    HTTP    `yaml:"http"`
		Swagger Swagger `yaml:"swagger"`
	}

	App struct {
		Name string `yaml:"name" env:"APP_NAME"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" env-required:"true"`
	}

	PG struct {
		Url     string `yaml:"url" env:"LOG_LEVEL" env-required:"true"`
		PoolMax int    `yaml:"pool_max" env:"PG_POOL_MAX" env-required:"true"`

		//host     string `yaml:"host" env:"PG_HOST" env-default:"localhost"`
		//port     int    `yaml:"port" env:"PG_PORT" env-default:"5432"`
		//user     string `yaml:"user" env:"PG_USER" env-required:"true"`
		//password string `yaml:"password" env:"PG_PASSWORD" env-required:"true"`
		//dbName   string `yaml:"db_name" env:"PG_DB_NAME"`
	}

	HTTP struct {
		Port        string        `yaml:"port" env:"HTTP_PORT" env-required:"true"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	}

	Swagger struct {
		Enabled *bool `yaml:"enabled" env:"SWAGGER_ENABLED" env-default:"false"`
	}
)

func NewConfig() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		log.Fatalf("config path is empty")
	}

	return NewConfigByPath(configPath)
}

func NewConfigByPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist")
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("failed to read config: " + err.Error())
	}

	return cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

package config

import (
	"flag"
	"fmt"
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
		JWT     JWT     `yaml:"jwt"`
	}

	App struct {
		Name string `yaml:"name" env:"APP_NAME"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" env-required:"true"`
	}

	PG struct {
		Url     string `yaml:"url"      env:"PG_URL"      env-required:"true"`
		PoolMax int    `yaml:"pool_max" env:"PG_POOL_MAX" env-required:"true"`

		//host     string `yaml:"host" env:"PG_HOST" env-default:"localhost"`
		//port     int    `yaml:"port" env:"PG_PORT" env-default:"5432"`
		//user     string `yaml:"user" env:"PG_USER" env-required:"true"`
		//password string `yaml:"password" env:"PG_PASSWORD" env-required:"true"`
		//dbName   string `yaml:"db_name" env:"PG_DB_NAME"`
	}

	HTTP struct {
		Port        string        `yaml:"port"         env:"HTTP_PORT" env-required:"true"`
		Timeout     time.Duration `yaml:"timeout"      env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	}

	Swagger struct {
		Enabled *bool `yaml:"enabled" env:"SWAGGER_ENABLED" env-default:"false"`
	}

	JWT struct {
		SignKey  string        `yaml:"sign_key"  env:"JWT_SIGN_KEY"  env-required:"true"`
		TokenTTL time.Duration `yaml:"token_ttl" env:"JWT_TOKEN_TTL" env-required:"true"`
	}
)

func NewConfig() *Config {
	path := fetchConfigPath()
	fmt.Println(path)
	if path == "" {
		log.Fatalf("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file does not exist")
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
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

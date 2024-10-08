package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	APIPort        string `envconfig:"API_PORT" required:"true"`
	PGDSN          string `envconfig:"PG_DSN" required:"true"`
	MaxHeaderBytes int    `envconfig:"MAX_HEADER_BYTES" required:"true"`
	ReadTimeout    int    `envconfig:"READ_TIMEOUT" required:"true"`
	WriteTimeout   int    `envconfig:"WRITE_TIMEOUT" required:"true"`
	LOG_LEVEL      string `envconfig:"LOG_LEVEL" required:"true"`
}

var Values config

func LoadFromFile(fpath string) error {
	godotenv.Load(fpath)

	err := envconfig.Process("", &Values)
	if err != nil {
		log.Printf("envconfig.Process(): %v", err.Error())
	}

	return err
}

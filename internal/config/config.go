package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Specification struct {
	SMTPPassword string `envconfig:"SMTP_PASSWORD"`
	SMTPHost     string `envconfig:"SMPT_HOST"`
	SMTPPort     string `envconfig:"SMTP_PORT"`
}

var Spec Specification

func InitConfig() {
	err := envconfig.Process("blaze", &Spec)
	if err != nil {
		log.Fatal(err.Error())
	}
}

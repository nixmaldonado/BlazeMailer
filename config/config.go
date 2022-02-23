package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Specification struct {
	Origin       string `envconfig:"from"`
	SMTPPassword string `envconfig:"SMTP_PASSWORD"`
	SMTPHost     string `envconfig:"SMPT_HOST" default:"smtp.gmail.com"`
	SMTPPort     string `envconfig:"SMTP_PORT" default:"587"`
}

var Spec Specification

func InitConfig() {
	err := envconfig.Process("blaze", &Spec)
	if err != nil {
		log.Fatal(err.Error())
	}
}

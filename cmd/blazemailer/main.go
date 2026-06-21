package main

import (
	"log"
	"net/http"

	"github.com/nixmaldonado/blazeMailer/internal/api"
	"github.com/nixmaldonado/blazeMailer/internal/config"
	"github.com/nixmaldonado/blazeMailer/internal/delivery"
)

func main() {
	config.InitConfig()
	sender := delivery.SMTPSender{
		Host:     config.Spec.SMTPHost,
		Port:     config.Spec.SMTPPort,
		Password: config.Spec.SMTPPassword,
	}
	srv := api.NewServer(sender)
	log.Println("listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", srv.Routes()))
}

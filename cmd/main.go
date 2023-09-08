package main

import (
	"companies/internal/app"
	"companies/internal/config"
	"log"
)

func main() {
	conf, err := config.InitConfig()
	if err != nil {
		log.Println("Config init failed")

		return
	}

	app.Run(conf)
}

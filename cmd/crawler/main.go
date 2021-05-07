package main

import (
	"log"

	"github.com/veotani/goblinator/pkg/config"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Print(err)
	} else {
		log.Print("Success!")
		log.Printf("Client ID: %s\n", config.BlizzardClientId)
		log.Printf("Client Secret: %s\n", config.BlizzardClientSecret)
	}
}

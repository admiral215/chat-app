package main

import (
	"chat-app/cmd"
	"chat-app/config"
	"log"
)

// @title Chat-App API
// @version 1.0
// @description This is a sample API Chat App using Gin and Swaggo
// @contact.name @rullyadmiral
// @contact.url https://github.com/admiral215
// @contact.email rullyadmiral@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		log.Fatal(err)
	}

	app, err := cmd.InitializeApp(cfg)
	if err != nil {
		log.Fatalf("cannot initialize app: %v", err)
	}

	err = app.Start()
	if err != nil {
		log.Fatalf("cannot start app: %v", err)
	}
}

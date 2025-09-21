package main

import (
	"log"
	"xeed/apps/cp-api/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

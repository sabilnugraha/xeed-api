package main

import (
	"log"
	"xeed/apps/cp-api/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

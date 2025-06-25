package main

import (
	"full-project-mock/internal/app"
	"log"
)

func main() {
	if err := app.RunApp(); err != nil {
		// TODO: посмотреть как лучше сделать тут
		log.Fatal(err)
	}
}

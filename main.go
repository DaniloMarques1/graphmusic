package main

import (
	"log"

	"github.com/danilomarques1/graphmusic/api"
	"github.com/joho/godotenv"
)

func main() {
	// TODO maybe should not log Fatal
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	s, err := api.NewServer("8080")
	if err != nil {
		log.Fatal(err)
	}
	s.Init()  // set things up
	s.Start() // start the http server
}

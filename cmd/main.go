package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/luxarts/go-oauth/internal/router"
	"log"
)

func main() {
	r := router.New()

	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}

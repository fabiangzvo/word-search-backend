package main

import (
	env "github.com/joho/godotenv"

	"word-search/api"
)

func main() {
	env.Load()

	api.InitServer()
}

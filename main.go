package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

type Config struct {
	API_KEY string
}

func main() {

	conf := Config{
		API_KEY: os.Getenv("API_KEY"),
	}

	ctx := context.Background()

	setupYoutubeAPI(ctx, conf, "youtube")

	createServer()
}

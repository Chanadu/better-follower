package main

import (
	"context"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	API_KEY string
}

func main() {

	conf := Config{
		API_KEY: os.Getenv("API_KEY"),
	}

	ctx := context.Background()

	yAPI, err := NewYoutubeAPI(ctx, conf)

	if err != nil {
		panic(err)
	}

	err = createServer(yAPI)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"fmt"
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

	details, recent, err := setupYoutubeAPI(ctx, conf, "jcole")
	if err != nil {
		panic(err)
	}
	fmt.Println(details)
	fmt.Println("Most Recent Video Title:", recent)

	createServer()
}

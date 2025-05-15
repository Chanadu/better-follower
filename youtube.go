package main

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func setupYoutubeAPI(ctx context.Context, conf Config, channelName string) {
	apiKey := string(conf.API_KEY)
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating YouTube service: %v", err)
	}

	// Search for the channel by name
	searchCall := youtubeService.Search.List([]string{"snippet"}).Q(channelName).Type("channel").MaxResults(1)
	searchResponse, err := searchCall.Do()
	if err != nil {
		log.Fatalf("Error searching for channel: %v", err)
	}

	if len(searchResponse.Items) == 0 {
		log.Fatalf("No channel found with the name: %s", channelName)
	}

	channelID := searchResponse.Items[0].Snippet.ChannelId

	// Retrieve channel details
	channelCall := youtubeService.Channels.List([]string{"snippet"}).Id(channelID)
	channelResponse, err := channelCall.Do()
	if err != nil {
		log.Fatalf("Error retrieving channel details: %v", err)
	}

	if len(channelResponse.Items) == 0 {
		log.Fatalf("No details found for channel ID: %s", channelID)
	}

	creationTime := channelResponse.Items[0].Snippet.PublishedAt
	log.Printf("Channel '%s' was created on: %s", channelName, creationTime)
}

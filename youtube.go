package main

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeAPI struct {
	API_KEY        string
	YoutubeService *youtube.Service
}

func NewYoutubeAPI(ctx context.Context, conf Config) (api *YoutubeAPI, err error) {
	apiKey := string(conf.API_KEY)
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube service: %w", err)
	}
	return &YoutubeAPI{
		apiKey,
		youtubeService
	}, nil
}

func setupYoutubeAPI(ctx context.Context, conf Config) (*YoutubeAPI, error) {
	apiKey := string(conf.API_KEY)
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		return "", "", fmt.Errorf("error creating YouTube service: %w", err)
	}

	// Search for channels by name
	// Search for the channel using its @tag
	searchCall := youtubeService.Channels.List([]string{"snippet"}).ForHandle(channelName)
	searchResponse, err := searchCall.Do()
	if err != nil {
		return "", "", fmt.Errorf("error searching for channels: %w", err)
	}

	if len(searchResponse.Items) == 0 {
		return "", "", fmt.Errorf("no channels found with the name: %s", channelName)
	}

	// Match the exact channel @tag
	channelID := searchResponse.Items[0].Id

	// Retrieve the most recent video using the Search endpoint
	videoSearchCall := youtubeService.Search.List([]string{"snippet"}).ChannelId(channelID).Type("video").Order("date").MaxResults(1)
	videoSearchResponse, err := videoSearchCall.Do()

	if err != nil {
		return "", "", fmt.Errorf("error retrieving most recent video: %w", err)
	}

	if len(videoSearchResponse.Items) == 0 {
		return "", "", fmt.Errorf("no videos found for channel ID: %s", channelID)
	}

	mostRecentVideoTitle := videoSearchResponse.Items[0].Snippet.Title

	// Retrieve channel details
	channelCall := youtubeService.Channels.List([]string{"snippet", "statistics"}).Id(channelID)
	channelResponse, err := channelCall.Do()
	if err != nil {
		return "", "", fmt.Errorf("error retrieving channel details: %w", err)
	}

	if len(channelResponse.Items) == 0 {
		return "", "", fmt.Errorf("no details found for channel ID: %s", channelID)
	}

	channel := channelResponse.Items[0]
	channelDetails := fmt.Sprintf("Channel Found:\nTitle: %s\nDescription: %s\nSubscribers: %d\nCreation Date: %s",
		channel.Snippet.Title, channel.Snippet.Description, channel.Statistics.SubscriberCount, channel.Snippet.PublishedAt)

	return channelDetails, mostRecentVideoTitle, nil
}

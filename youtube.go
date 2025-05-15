package main

import (
	"context"
	"fmt"
	"net/http"

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
		youtubeService,
	}, nil
}

func (yAPI *YoutubeAPI) SearchByHandle(handle string) (channelDetails string, mostRecentVideoTitle string, err error) {

	searchCall := yAPI.YoutubeService.Channels.List([]string{"snippet"}).ForHandle(handle)
	searchResponse, err := searchCall.Do()
	if err != nil {
		return "", "", fmt.Errorf("error searching for channels: %w", err)
	}

	if len(searchResponse.Items) == 0 {
		return "", "", fmt.Errorf("no channels found with the name: %s", handle)
	}

	channelID := searchResponse.Items[0].Id

	videoSearchCall := yAPI.YoutubeService.Search.List([]string{"snippet"}).
		ChannelId(channelID).
		Type("video").
		Order("date").
		MaxResults(1)

	videoSearchResponse, err := videoSearchCall.Do()

	if err != nil {
		return "", "", fmt.Errorf("error retrieving most recent video: %w", err)
	}

	if len(videoSearchResponse.Items) == 0 {
		return "", "", fmt.Errorf("no videos found for channel: %s (ID: %s)", handle, channelID)
	}

	mostRecentVideoTitle = videoSearchResponse.Items[0].Snippet.Title

	channelCall := yAPI.YoutubeService.Channels.List([]string{"snippet", "statistics"}).Id(channelID)
	channelResponse, err := channelCall.Do()
	if err != nil {
		return "", "", fmt.Errorf("error retrieving channel details: %w", err)
	}

	if len(channelResponse.Items) == 0 {
		return "", "", fmt.Errorf("no details found for channel ID: %s", channelID)
	}

	channel := channelResponse.Items[0]
	channelDetails = fmt.Sprintf("Channel Found:\nTitle: %s\nDescription: %s\nSubscribers: %d\nCreation Date: %s",
		channel.Snippet.Title, channel.Snippet.Description, channel.Statistics.SubscriberCount, channel.Snippet.PublishedAt)

	return channelDetails, mostRecentVideoTitle, nil
}

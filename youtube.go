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
		youtubeService,
	}, nil
}

type ChannelData struct {
	Id              string
	Name            string
	Description     string
	SubscriberCount uint64
	CreationTime    string
}

func (c *ChannelData) ToString() string {
	return fmt.Sprintf("Id: %s\nName: %s\nDescription: %s\nSubscriber Count: %d\nCreation Time: %s",
		c.Id, c.Name, c.Description, c.SubscriberCount, c.CreationTime)
}

type VideoData struct {
	Title        string
	Description  string
	CreationTime string
}

func (v *VideoData) ToString() string {
	return fmt.Sprintf("Title: %s\nDescription: %s\nCreation Time: %s",
		v.Title, v.Description, v.CreationTime)
}

func (yAPI *YoutubeAPI) RecentVideosByChannelId(channelID string, count int) (videoData []*VideoData, err error) {
	if count <= 0 {
		return videoData, fmt.Errorf("count must be greater than 0")
	}
	if count > 100 {
		return videoData, fmt.Errorf("count must be less than 101")
	}

	videoSearchCall := yAPI.YoutubeService.Search.List([]string{"snippet"}).
		ChannelId(channelID).
		Type("video").
		Order("date").
		MaxResults(int64(count))

	videoSearchResponse, err := videoSearchCall.Do()

	if err != nil {
		return videoData, fmt.Errorf("error retrieving most recent video(s): %w", err)
	}

	if len(videoSearchResponse.Items) == 0 {
		return videoData, fmt.Errorf("no videos found for channel id: %s", channelID)
	}

	videoData = make([]*VideoData, len(videoSearchResponse.Items))
	for i, item := range videoSearchResponse.Items {
		vd := &VideoData{
			Title:        item.Snippet.Title,
			CreationTime: item.Snippet.PublishedAt,
			Description:  item.Snippet.Description,
		}
		videoData[i] = vd
	}

	return videoData, nil
}

func (yAPI *YoutubeAPI) SearchByHandle(handle string) (channelData *ChannelData, err error) {
	searchCall := yAPI.YoutubeService.Channels.List([]string{"snippet"}).ForHandle(handle)
	searchResponse, err := searchCall.Do()
	if err != nil {
		return channelData, fmt.Errorf("error searching for channels: %w", err)
	}

	if len(searchResponse.Items) == 0 {
		return channelData, fmt.Errorf("no channels found with the name: %s", handle)
	}

	channelID := searchResponse.Items[0].Id
	channelCall := yAPI.YoutubeService.Channels.List([]string{"snippet", "statistics"}).Id(channelID)
	channelResponse, err := channelCall.Do()
	if err != nil {
		return channelData, fmt.Errorf("error retrieving channel details: %w", err)
	}

	if len(channelResponse.Items) == 0 {
		return channelData, fmt.Errorf("no details found for channel ID: %s", channelID)
	}

	channel := channelResponse.Items[0]
	channelData = &ChannelData{
		Id:              channelID,
		Name:            channel.Snippet.Title,
		Description:     channel.Snippet.Description,
		SubscriberCount: channel.Statistics.SubscriberCount,
		CreationTime:    channel.Snippet.PublishedAt,
	}

	return channelData, nil
}

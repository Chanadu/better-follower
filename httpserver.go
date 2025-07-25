package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "home\n")
}

func (yAPI *YoutubeAPI) RedirectToHandleSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, req.URL.Path+"/1", http.StatusFound)
	}
}

func (yAPI *YoutubeAPI) HandleSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		args := strings.Split(req.URL.Path, "/")
		handle := args[2]
		strCount := args[3]

		if strCount == "" {
			panic("count is empty")
		}
		count, err := strconv.Atoi(strCount)
		if err != nil {
			http.Error(w, "Invalid count parameter (Cannot convert to integer)", http.StatusBadRequest)
			return
		}

		channelData, err := yAPI.SearchByHandle(handle)
		mostRecentVideoData, err := yAPI.RecentVideosByChannelId(channelData.Id, count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Channel Details: %s", channelData.ToString())
		for _, video := range mostRecentVideoData {
			fmt.Fprintf(w, "Video Details: %s", video.ToString())
		}
	}
}

func createServer(yAPI *YoutubeAPI) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/artist/{handle}", yAPI.RedirectToHandleSearchHandler())
	mux.HandleFunc("/artist/{handle}/{count}", yAPI.HandleSearchHandler())

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

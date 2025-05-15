package main

import (
	"fmt"
	"net/http"
)

func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "home\n")
}

func (yAPI *YoutubeAPI) HandleSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handle := req.URL.Path[len("/artist/"):]
		channelDetails, mostRecentVideoTitle, err := yAPI.SearchByHandle(handle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Channel Details: %s\nMost Recent Video Title: %s\n", channelDetails, mostRecentVideoTitle)
	}
}

func createServer(yAPI *YoutubeAPI) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/artist/{handle}", yAPI.HandleSearchHandler())

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

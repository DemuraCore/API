package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func getNowPlaying(c *gin.Context) {
	accessToken, err := getAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	req, err := http.NewRequest("GET", nowPlayingEndpoint, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch now playing"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		getRecentlyPlayedHelper(c)
		return
	}

	var nowPlaying map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&nowPlaying)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode now playing response"})
		return
	}

	c.JSON(http.StatusOK, nowPlaying)
}

func getRecentlyPlayedHelper(c *gin.Context) {
	accessToken, err := getAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/recently-played", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recently played tracks"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recently played tracks"})
		return
	}

	var data struct {
		Items []struct {
			Track struct {
				Name    string `json:"name"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
				Album struct {
					Name   string `json:"name"`
					Images []struct {
						URL string `json:"url"`
					} `json:"images"`
				} `json:"album"`
				ExternalURLs struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
			} `json:"track"`
		} `json:"items"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode recently played tracks response"})
		return
	}

	if len(data.Items) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No recently played tracks found"})
		return
	}

	recentTrack := data.Items[0]
	trackName := recentTrack.Track.Name
	artistName := recentTrack.Track.Artists[0].Name
	albumName := recentTrack.Track.Album.Name
	trackPhoto := recentTrack.Track.Album.Images[0].URL
	songUrl := recentTrack.Track.ExternalURLs.Spotify

	c.JSON(http.StatusOK, gin.H{
		"isPlaying":  false,
		"trackName":  trackName,
		"artistName": artistName,
		"albumName":  albumName,
		"trackPhoto": trackPhoto,
		"songUrl":    songUrl,
	})
}

func getRecentlyPlayed(c *gin.Context) {
	getRecentlyPlayedHelper(c)
}
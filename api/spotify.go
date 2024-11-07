package api

import (
	"DemuraCore/API/model"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	clientID     = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	refreshToken = os.Getenv("SPOTIFY_REFRESH_TOKEN")

	tokenEndpoint          = "https://accounts.spotify.com/api/token"
	nowPlayingEndpoint     = "https://api.spotify.com/v1/me/player/currently-playing"
	recentlyPlayedEndpoint = "https://api.spotify.com/v1/me/player/recently-played"
	basicAuth              = base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
)

func getAccessToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+basicAuth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResponse model.AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func makeSpotifyRequest(endpoint, accessToken string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func GetNowPlaying(c *gin.Context) {
	accessToken, err := getAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	resp, err := makeSpotifyRequest(nowPlayingEndpoint, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch now playing"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		GetRecentlyPlayed(c)
		return
	}

	var nowPlaying map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&nowPlaying); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode now playing response"})
		return
	}

	c.JSON(http.StatusOK, nowPlaying)
}

func GetRecentlyPlayed(c *gin.Context) {
	accessToken, err := getAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	resp, err := makeSpotifyRequest(recentlyPlayedEndpoint, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recently played tracks"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recently played tracks"})
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(&model.DataRecentPlaying); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode recently played tracks response"})
		return
	}

	if len(model.DataRecentPlaying.Items) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No recently played tracks found"})
		return
	}

	recentTrack := model.DataRecentPlaying.Items[0]
	c.JSON(http.StatusOK, gin.H{
		"isPlaying":  false,
		"trackName":  recentTrack.Track.Name,
		"artistName": recentTrack.Track.Artists[0].Name,
		"albumName":  recentTrack.Track.Album.Name,
		"trackPhoto": recentTrack.Track.Album.Images[0].URL,
		"songUrl":    recentTrack.Track.ExternalURLs.Spotify,
	})
}

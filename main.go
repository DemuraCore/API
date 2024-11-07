package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	clientID     = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	refreshToken = os.Getenv("SPOTIFY_REFRESH_TOKEN")
)

func main() {
	log.Println("Starting server...")
	log.Println("Listening on localhost:3000")

	log.Println("Client ID:", clientID)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/now-playing", getNowPlaying)
	router.GET("/recently-played", getRecentlyPlayed)

	router.Run("0.0.0.0:3000")
}

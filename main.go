package main

import (
	"DemuraCore/API/api"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting DemuraCore API V2.0...")
	log.Println("Listening on PORT 3000")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(api.AuthMiddleware())

	// Spotify api
	router.GET("/util/spotify/now-playing", api.GetNowPlaying)
	router.GET("/util/spotify/recently-played", api.GetRecentlyPlayed)

	// Heartbeat api
	router.GET("/heartbeat", api.Heartbeat)

	// Discord api
	router.GET("/core/discord/me", api.GetMeDetails)

	router.Run("0.0.0.0:3000")
}

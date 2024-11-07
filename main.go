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
	Router := gin.Default()

	/*
		Core API DemuraCore API V2.0
		- Requires Authorization
	*/

	CoreRouter := Router.Group("/core")
	CoreRouter.Use(api.AuthMiddleware())

	CoreRouter.GET("/discord/me", api.GetMeDetails)

	/*
		Util API DemuraCore API V2.0
		- No Authorization required
		- Public API
	*/
	UtilityRouter := Router.Group("/util")

	// Spotify api
	UtilityRouter.GET("/spotify/now-playing", api.GetNowPlaying)
	UtilityRouter.GET("/spotify/recently-played", api.GetRecentlyPlayed)

	// Heartbeat api
	Router.GET("/heartbeat", api.Heartbeat)

	Router.Run("0.0.0.0:3000")
}

package api

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	// Discord API
	URL   = "https://discord.com/api/v10"
	TOKEN = os.Getenv("DISCORD_TOKEN")
)

func GetMeDetails(c *gin.Context) {
	url := "https://discord.com/api/v10/users/@me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Authorization", TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make request"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}

// func SendMessage(c *gin.Context) {
// 	url := "https://discord.com/api/v10/channels/1234567890/messages"

// 	req, err := http.NewRequest("POST", url, nil)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
// 		return
// 	}

// 	req.Header.Set("Authorization", TOKEN)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make request"})
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
// 		return
// 	}

// 	c.Data(resp.StatusCode, "application/json", body)
// }

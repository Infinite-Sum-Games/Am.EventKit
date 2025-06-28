package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	"github.com/gin-gonic/gin"
)

func main() {
	// setting up application configuration
	config, err := cmd.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cmd.Env = config
	log.Println("Environment variables loaded successfully.")

	r := gin.New()
	r.Use(gin.Logger())
	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is live ◪_◪",
		})
	})
	err = r.Run(":" + "9000")
	if err != nil {
		fmt.Println("Server failed")
		return
	}
}

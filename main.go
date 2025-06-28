package main

import (
	"fmt"
	"net/http"

	services "github.com/Thanus-Kumaar/anokha-2025-backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	err := services.InitDBPool()
	if err != nil {
		panic(fmt.Errorf("Failed to initialize database pool: %w", err))
	}

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

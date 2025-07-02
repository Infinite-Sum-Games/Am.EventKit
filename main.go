package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	"github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// initializing the logger and other middlewares
	logger, err := cmd.InitLogger("DEVELOPMENT") // NOTE: hardcoded now, will change once viper is setup
	if err != nil {
		log.Fatalf("Logger initialization failed: %v", err)
	}
	pkg.Log = logger
	pkg.Log.LogInfo("Logger initiation successful")

	r.Use(middleware.RequestLoggerMiddleware(pkg.Log))

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

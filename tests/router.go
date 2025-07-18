package tests

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var mockEnv = struct {
	Domain string
}{
	Domain: "http://example.com",
}

func SetupTestRouter() *gin.Engine {

	gin.SetMode(gin.TestMode)
	r := gin.New()
	config := cors.Config{
		AllowOrigins:              []string{mockEnv.Domain},
		AllowWildcard:             true,
		AllowMethods:              []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:              []string{"X-Csrf-Token", "Origin", "Content-Type"},
		AllowCredentials:          true,
		OptionsResponseStatusCode: 204,
		MaxAge:                    12 * time.Hour,
	}
	r.Use(cors.New(config))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is live ◪_◪",
		})
	})

	return r
}

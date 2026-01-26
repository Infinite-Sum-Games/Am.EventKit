package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupMCPRouter() *gin.Engine {
	config := cors.Config{
		AllowOrigins:              []string{"*"},
		AllowWildcard:             true,
		AllowMethods:              []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:              []string{"X-Csrf-Token", "Origin", "Content-Type"},
		AllowCredentials:          true,
		OptionsResponseStatusCode: 204,
		MaxAge:                    12 * time.Hour,
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(cors.New(config))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "MCP server is running",
			"port":   Env.App.MCPPort,
		})
	})

	return r
}

func RunMCPServer(ctx context.Context) error {
	if Env == nil {
		config, err := LoadConfig()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
		Env = config
	}

	router := SetupMCPRouter()

	port := Env.App.MCPPort
	if port == 0 {
		port = 8081
	}

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router,
	}

	go func() {
		log.Printf("[OK]: Starting MCP server on port %d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("MCP server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("[OK]: Shutting down MCP server...")
	case <-ctx.Done():
		log.Println("[OK]: Context cancelled, shutting down MCP server...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("MCP server forced to shutdown: %v", err)
		return err
	}

	log.Println("[OK]: MCP server shutdown complete")
	return nil
}

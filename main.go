package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	"github.com/Thanus-Kumaar/anokha-2025-backend/internal/auth"
	"github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	// setting up application configuration
	config, err := cmd.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cmd.Env = config
	log.Println("environment variables loaded successfully")

	// initialize RSA keys if missing and verify them
	if err := auth.InitRSAKeys(config.PrivateKeyPath, config.PublicKeyPath); err != nil {
		log.Fatalf("failed to initialize RSA keys: %v", err)
	}

	if err := auth.VerifyRSAKeys(config.PrivateKeyPath, config.PublicKeyPath); err != nil {
		log.Fatalf("rsa key verification failed: %v", err)
	}

	// load keys into memory
	privKey, pubKey, err := auth.LoadRSAKeysFromFiles(config.PrivateKeyPath, config.PublicKeyPath)
	if err != nil {
		log.Fatalf("failed to load rsa keys into memory: %v", err)
	}
	auth.SetRSAKeys(privKey, pubKey)
	log.Println("rsa keys loaded and stored in memory")

	// ensure logs folder exists
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatalf("failed to create logs directory: %v", err)
	}

	// initialize the database pool
	if err := cmd.InitDBPool(); err != nil {
		panic(fmt.Errorf("failed to initialize database pool: %w", err))
	}

	// initialize logger and middlewares
	logger, err := cmd.InitLogger("DEVELOPMENT") // later replace hardcoded env
	if err != nil {
		log.Fatalf("logger initialization failed: %v", err)
	}
	pkg.Log = logger
	pkg.Log.LogInfo("logger initiation successful")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.RequestLoggerMiddleware(pkg.Log))

	// sample test route
	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "server is live ◪_◪",
		})
	})

	// start server
	if err := r.Run(":9000"); err != nil {
		fmt.Println("server failed")
	}
}

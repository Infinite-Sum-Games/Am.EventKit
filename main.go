package main

import (
	"fmt" "log" "net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	// setting up application configuration
	config, err := cmd.LoadConfig()
	if err != nil {
		log.Fatalf("%v", err)
	}
	cmd.Env = config
	log.Println("Environment variables loaded successfully.")

	// initializing the logger and other middlewares
	logger, err := cmd.InitLogger("DEVELOPMENT") // NOTE: hardcoded now, will change once viper is setup
	if err != nil {
		log.Fatalf("Logger initialization failed: %v", err)
	}
	pkg.Log = logger
	pkg.Log.LogInfo("Logger initiation successful")

	// Initialize RSA
	err = cmd.CheckRSAKeyPairExists()
	if err != nil {
		err = cmd.GenerateRSAKeyPair()
		if err != nil {
			panic(fmt.Errorf("Failed to initialize RSA: %w", err))
		}
		pkg.Log.LogInfo("[OK]: RSA keypair generated and saved successfully.")
	} else {
		pkg.Log.LogInfo("[OK]: Using existing RSA keypair.")
	}

	err = cmd.InitDBPool()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database pool: %w", err))
	}


	r.Use(pkg.Log.LogRequest)
	r := gin.New()
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

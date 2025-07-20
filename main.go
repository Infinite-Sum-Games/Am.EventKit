package main

import (
	"log"
	"net/http"
	"time"

	apiAuth "github.com/Thanus-Kumaar/anokha-2025-backend/api/auth"
	apiEvent "github.com/Thanus-Kumaar/anokha-2025-backend/api/event"
	apiProfile "github.com/Thanus-Kumaar/anokha-2025-backend/api/profile"
	apiStaff "github.com/Thanus-Kumaar/anokha-2025-backend/api/staff"
	apiTag "github.com/Thanus-Kumaar/anokha-2025-backend/api/tag"
	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	config := cors.Config{
		AllowOrigins:              []string{cmd.Env.Domain},
		AllowWildcard:             true,
		AllowMethods:              []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:              []string{"X-Csrf-Token", "Origin", "Content-Type"},
		AllowCredentials:          true,
		OptionsResponseStatusCode: 204,
		MaxAge:                    12 * time.Hour,
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(cors.New(config)) // Setup CORS() first before other middlewares
	r.Use(pkg.Log.LogMiddleware)
	r.Use(pkg.TagRequestWithId)
	r.Use(mw.RecoveryPanics)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is live ◪_◪",
		})
		pkg.Log.SuccessCtx(c)
	})

	v1 := r.Group("/api/v1")
	authRouter := v1.Group("/auth")
	staffRouter := v1.Group("/staff")
	userRouter := v1.Group("/user")

	apiAuth.StudentAuthRoutes(authRouter)
	apiAuth.StaffAuthRoutes(authRouter)
	apiProfile.ProfileRoutes(userRouter)
	apiEvent.EventRoutes(userRouter)
	apiStaff.AttendanceRoutes(staffRouter)
	apiTag.TagRoutes(v1)

	return r
}

func StartApp() {
	// Setting up environment variables
	config, err := cmd.LoadConfig()
	if err != nil {
		log.Printf("[CRASH] Failed to load environment variables: %v", err)
		return
	}
	cmd.Env = config
	log.Println("[OK]: Environment variables loaded successfully.")

	// Initializing the logger and other middlewares
	pkg.Log, err = pkg.InitLogger(cmd.Env.Environment)
	if err != nil {
		log.Printf("[CRASH]: Logger initialization failed: %v", err)
		return
	}
	pkg.Log.Info("[OK]: Logger initiation successful")

	// Initialize RSA
	err = cmd.CheckRSAKeyPairExists()
	if err != nil {
		err = cmd.GenerateRSAKeyPair()
		if err != nil {
			pkg.Log.Fatal("[CRASH]: Failed to initialize rsa", err)
		}
		pkg.Log.Info("[OK]: RSA keypair generated and saved successfully.")
	} else {
		pkg.Log.Info("[OK]: Using existing RSA keypair.")
	}

	// Setup PASETO
	if err := pkg.InitPaseto(); err != nil {
		pkg.Log.Fatal("[CRASH]: Paseto initialization failed", err)
	}
	pkg.Log.Info("[OK]: Paseto initialization successful!")

	// Initialize DB Pool
	err = cmd.InitDBPool()
	if err != nil {
		pkg.Log.Fatal("[CRASH]: Failed to initialize database pool", err)
		return
	}
	pkg.Log.Info("[OK]: Initialized database pool successfully")

	// Initialize Valkey (cache)
	cmd.Valkey, err = cmd.InitValkey()
	if err != nil {
		pkg.Log.Fatal("[CRASH]: Failed to initialize cache", err)
		return
	}
	pkg.Log.Info("[OK]: Valkey initialized successfully")

	// Initialize server
	pkg.Log.Info("[OK]: Start the server on port 9000")
	err = SetupRouter().Run(":" + "9000")
	if err != nil {
		pkg.Log.Fatal("[CRASH]: Server failed to start", err)
		return
	}
}

func main() {
	StartApp()
}

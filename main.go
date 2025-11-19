package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiAttend "github.com/Thanus-Kumaar/anokha-2025-backend/api/attendance"
	apiAuth "github.com/Thanus-Kumaar/anokha-2025-backend/api/auth"
	apiEvent "github.com/Thanus-Kumaar/anokha-2025-backend/api/event"
	apiMail "github.com/Thanus-Kumaar/anokha-2025-backend/api/mail"
	apiOrganizer "github.com/Thanus-Kumaar/anokha-2025-backend/api/organizers"
	apiProfile "github.com/Thanus-Kumaar/anokha-2025-backend/api/profile"
	apiTag "github.com/Thanus-Kumaar/anokha-2025-backend/api/tag"
	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	"github.com/Thanus-Kumaar/anokha-2025-backend/mail"
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(mailerSvc *mail.MailerService) *gin.Engine {

	config := cors.Config{
		AllowOrigins:              []string{cmd.Env.ClientDomain},
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

	r.GET("/test", mw.PrometheusMiddleware("test"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is live ◪_◪",
		})
		pkg.Log.SuccessCtx(c)
	})

	mailController := apiMail.NewController(mailerSvc)
	apiMail.SetRoutes(r, mailController)

	v1 := r.Group("/api/v1")
	authRouter := v1.Group("/auth")
	attendanceRouter := v1.Group("/attendance")
	userRouter := v1.Group("/user")
	eventRouter := v1.Group("/events")

	apiAuth.StudentAuthRoutes(authRouter)
	apiAuth.OrganizerAuthRoutes(authRouter)
	apiProfile.ProfileRoutes(userRouter)
	apiEvent.EventRoutes(eventRouter)
	apiAttend.AttendanceRoutes(attendanceRouter)
	apiTag.TagRoutes(userRouter)
	apiOrganizer.OrganizerRoutes(userRouter)

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

	// Initialize Mailer Service
	mailerSvc, err := mail.NewMailerService("mail/mail-queue", 4)
	if err != nil {
		pkg.Log.Fatal("failed to create mailer service", err)
	}
	mailerSvc.Start()
	pkg.Log.Info("[OK]: Mailer service started successfully")

	// Initialize server
	server := &http.Server{
		Addr:    ":" + "9000",
		Handler: SetupRouter(mailerSvc),
	}

	go func() {
		pkg.Log.Info("[OK]: Start the server on port 9000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pkg.Log.Fatal("could not listen on port 9000", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	pkg.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		pkg.Log.Fatal("Server forced to shutdown", err)
	}

	mailerSvc.Shutdown()

	pkg.Log.Info("Server exiting")
}

func main() {
	StartApp()
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	apiAnalytics "github.com/Thanus-Kumaar/anokha-2025-backend/api/analytics"
	apiAttend "github.com/Thanus-Kumaar/anokha-2025-backend/api/attendance"
	apiAuth "github.com/Thanus-Kumaar/anokha-2025-backend/api/auth"
	apiBooking "github.com/Thanus-Kumaar/anokha-2025-backend/api/booking"
	apiEvent "github.com/Thanus-Kumaar/anokha-2025-backend/api/event"
	apiOrganizers "github.com/Thanus-Kumaar/anokha-2025-backend/api/organizers"
	apiPeople "github.com/Thanus-Kumaar/anokha-2025-backend/api/people"
	apiProfile "github.com/Thanus-Kumaar/anokha-2025-backend/api/profile"
	apiTag "github.com/Thanus-Kumaar/anokha-2025-backend/api/tag"
	mq "github.com/Thanus-Kumaar/anokha-2025-backend/message-queue"

	cmd "github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	mail "github.com/Thanus-Kumaar/anokha-2025-backend/mail"
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	pkg "github.com/Thanus-Kumaar/anokha-2025-backend/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(mailerSvc *mail.MailerService) *gin.Engine {

	config := cors.Config{
		AllowOrigins:              []string{cmd.Env.ClientDomain},
		AllowWildcard:             true,
		AllowMethods:              []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:              []string{"X-Csrf-Token", "Origin", "Content-Type"},
		AllowCredentials:          true,
		OptionsResponseStatusCode: 204,
		MaxAge:                    12 * time.Hour,
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(cors.New(config)) // Setup CORS() first before other middlewares
	r.Use(mw.MaintainanceMiddleware)
	r.Use(pkg.Log.LogMiddleware)
	r.Use(pkg.TagRequestWithId)
	r.Use(mw.RecoveryPanics)
	r.Use(mw.PrometheusMiddleware("anokha-26"))

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is live ◪_◪",
		})
		pkg.Log.SuccessCtx(c)
	})

	r.GET("/metrics", mw.MetricsHandler())

	v1 := r.Group("/api/v1")
	authRouter := v1.Group("/auth")
	attendanceRouter := v1.Group("/attendance")
	userRouter := v1.Group("/user")
	eventRouter := v1.Group("/events")
	peopleRouter := v1.Group("/people")
	tagRouter := v1.Group("/tags")
	organizerRouter := v1.Group("/organizers")
	analyticsRouter := v1.Group("/analytics")

	apiAuth.StudentAuthRoutes(authRouter)
	apiAuth.OrganizerAuthRoutes(authRouter)
	apiAuth.AdminAuthRoutes(authRouter)
	apiProfile.ProfileRoutes(userRouter)
	apiEvent.EventRoutes(eventRouter)
	apiTag.TagRoutes(tagRouter)
	apiAttend.AttendanceRoutes(attendanceRouter)
	apiPeople.PeopleRoutes(peopleRouter)
	apiOrganizers.OrganizerRoutes(organizerRouter)
	apiOrganizers.OrganizerDashboardRoutes(organizerRouter)
	apiBooking.BookingRoutes(eventRouter)
	apiAnalytics.AnalyticsRoutes(analyticsRouter)

	return r
}

func StartApp() {
	// Setting up environment variables
	config, err := cmd.LoadConfig()
	if err != nil {
		log.Printf("[CRASH]: Failed to load environment variables: %v", err)
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

	// Initialize Redis for caching and rate limiing
	cmd.Redis, err = cmd.InitRedis()
	if err != nil {
		pkg.Log.Info("[CRASH]: Redis failed to start")
		return
	}
	pkg.Log.Info("[OK]: Redis serivce started successfully")

	// Initialize RabbitMQ for durable message passing
	mq.Rabbit, err = mq.NewBroker(cmd.Env.MsgBrokerConnUrl)
	if err != nil {
		return
	}
	pkg.Log.Info("[OK]: Message broker initialized successfully.")

	// Setup LIVE flags
	pkg.InitFlag()

	// Initialize Mailer Service
	mail.Mail, err = mail.NewMailerService("mail/mail-queue", 4)
	if err != nil {
		pkg.Log.Fatal("[CRASH]: failed to create mailer service", err)
		return
	}
	mail.Mail.Start()
	pkg.Log.Info("[OK]: Mailer service started successfully")

	// Initialize server
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cmd.Env.Port),
		Handler: SetupRouter(mail.Mail),
	}

	go func() {
		portStr := strconv.Itoa(cmd.Env.Port)
		pkg.Log.Info("[OK]: Start the server on port " + portStr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pkg.Log.Fatal("could not listen on port "+portStr, err)
		} // Blocking in nature (?)
	}()

	// Graceful shutdown with 10 second timeout; no new connections accepted
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	pkg.Log.Info("[OK]: Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		pkg.Log.Fatal("[OK]: Server forced to shutdown", err)
	}

	// Mailer shutdown sequence
	mail.Mail.Shutdown()

	pkg.Log.Info("[OK]: Server shutting down")
}

func main() {
	StartApp()
}

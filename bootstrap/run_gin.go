package bootstrap

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Infinite-Sum-Games/Am.EventKit/common"
	"github.com/Infinite-Sum-Games/Am.EventKit/configs"
	"github.com/Infinite-Sum-Games/Am.EventKit/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func initializeRouter(cfg *configs.Config) *gin.Engine {
	clientDomain := "http://localhost:3000"
	if cfg.App.ClientDomain != "" {
		clientDomain = cfg.App.ClientDomain
	}

	config := cors.Config{
		AllowOrigins:              []string{clientDomain},
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
	r.Use(common.TagRequestWithId)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is live ◪_◪",
		})
	})

	router.InitRouter(r, cfg)

	return r
}

func (a *App) RunGinServer() error {
	config := a.appConfig
	router := initializeRouter(a.appConfig)
	// msg := fmt.Sprintf("starting server on port - %d", config.App.Port)
	// TODO: Fix after better logger is introduced
	// logger.Log.Info(msg)

	srv := &http.Server{
		Addr:    strings.Join([]string{"0.0.0.0", strconv.Itoa(config.App.Port)}, ":"),
		Handler: router,
		// TODO: Introduce these variables
		// ReadTimeout: time.Duration(config.Http.ReadTimeout) * time.Second,
		// WriteTimeout: time.Duration(config.Http.WriteTimeout) * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// TODO: Fix after logger
			// logger.Log.Error("gin server errored out", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// logger.Log.Info("marking server as unhealthy")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// logger.Log.Error()
	}
	<-ctx.Done()
	return nil
}

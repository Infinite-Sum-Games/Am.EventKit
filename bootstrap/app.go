package bootstrap

import (
	"context"
	"log"

	"github.com/Infinite-Sum-Games/Am.EventKit/configs"
	"github.com/Infinite-Sum-Games/Am.EventKit/registry"
)

type App struct {
	registry  *registry.Registry
	appConfig *configs.Config
}

func (a *App) RunMCPServer(ctx context.Context) error {
	return RunMCPServer(ctx)
}

func (a *App) Close() {
	if a.registry == nil {
		return
	}

	// Close Redis if available
	if redisCli := a.registry.GetRedisClient(); redisCli != nil {
		if err := redisCli.Close(); err != nil {
			log.Printf("[ERROR]: Failed to close Redis client: %v", err)
		}
	}

	// Close OLTP pool if available
	if oltpPool := a.registry.GetOLTPPool(); oltpPool != nil {
		oltpPool.Close()
	}

	// Close OLAP pool if available
	if olapPool := a.registry.GetOLAPPool(); olapPool != nil {
		if closer, ok := olapPool.(interface{ Close() }); ok {
			closer.Close()
		}
	}

	log.Println("[OK]: App closed")
}

func NewApp() *App {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	configs.Env = config

	reg := registry.New(
		registry.WithConfig(config),
	)
	registry.SetSingletonObject(reg)

	if err := registry.InitializeServices(); err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	return &App{
		registry:  reg,
		appConfig: config,
	}
}

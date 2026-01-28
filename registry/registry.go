package registry

import (
	"sync"

	"github.com/Infinite-Sum-Games/Am.EventKit/configs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Registry struct {
	mu             sync.RWMutex
	kv             map[string]string
	config         configs.Config
	redisCli       *redis.Client
	oltpPool       *pgxpool.Pool
	olapPool       *pgxpool.Pool
	mailerService  interface{}
	paymentService interface{}
	messageQueue   interface{}
	// flagsmithClient *flagsmith.Client // Flagsmith feature flag client
}

var (
	registry   *Registry
	registryMu sync.RWMutex
)

func New(options ...func(r *Registry)) *Registry {
	reg := &Registry{
		kv: make(map[string]string),
	}
	for _, option := range options {
		option(reg)
	}
	return reg
}

func WithConfig(cfg *configs.Config) func(*Registry) {
	return func(reg *Registry) {
		reg.config = *cfg
	}
}

func WithRedisClient(cli *redis.Client) func(*Registry) {
	return func(reg *Registry) {
		reg.redisCli = cli
	}
}

func WithOLTPPool(pool *pgxpool.Pool) func(*Registry) {
	return func(reg *Registry) {
		reg.oltpPool = pool
	}
}

func WithOLAPPool(pool interface{}) func(*Registry) {
	return func(reg *Registry) {
		reg.olapPool = pool
	}
}

func WithMailerService(service interface{}) func(*Registry) {
	return func(reg *Registry) {
		reg.mailerService = service
	}
}

func WithPaymentService(service interface{}) func(*Registry) {
	return func(reg *Registry) {
		reg.paymentService = service
	}
}

func WithMessageQueue(queue interface{}) func(*Registry) {
	return func(reg *Registry) {
		reg.messageQueue = queue
	}
}

// func WithFlagsmithClient(client *flagsmith.Client) func(*Registry) {
// 	return func(reg *Registry) {
// 		reg.flagsmithClient = client
// 	}
// }

func WithKV(key, val string) func(*Registry) {
	return func(reg *Registry) {
		reg.kv[key] = val
	}
}

func (reg *Registry) GetConfig() configs.Config {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	return reg.config
}

func (reg *Registry) GetRedisClient() *redis.Client {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	return reg.redisCli
}

func (reg *Registry) GetOLTPPool() *pgxpool.Pool {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	return reg.oltpPool
}

func (reg *Registry) GetOLAPPool() interface{} {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	return reg.olapPool
}

func (reg *Registry) GetMailerService() interface{} {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	return reg.mailerService
}

func (reg *Registry) GetPaymentService() interface{} {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	return reg.paymentService
}

func (reg *Registry) GetMessageQueue() interface{} {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	return reg.messageQueue
}

// func (reg *Registry) GetFlagsmithClient() *flagsmith.Client {
// 	reg.mu.RLock()
// 	defer reg.mu.RUnlock()
// 	return reg.flagsmithClient
// }

func (reg *Registry) GetVal(key string) string {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	val, ok := reg.kv[key]
	if !ok {
		return ""
	}
	return val
}

func GetSingletonObject() *Registry {
	registryMu.RLock()
	defer registryMu.RUnlock()
	return registry
}

func SetSingletonObject(r *Registry) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = r
}

func GetConfig() configs.Config {
	r := GetSingletonObject()
	return r.config
}

func GetRedisClient() *redis.Client {
	r := GetSingletonObject()
	return r.redisCli
}

func GetOLTPPool() *pgxpool.Pool {
	r := GetSingletonObject()
	return r.oltpPool
}

func GetOLAPPool() interface{} {
	r := GetSingletonObject()
	return r.olapPool
}

func GetMailerService() interface{} {
	r := GetSingletonObject()
	return r.mailerService
}

func GetPaymentService() interface{} {
	r := GetSingletonObject()
	return r.paymentService
}

func GetMessageQueue() interface{} {
	r := GetSingletonObject()
	return r.messageQueue
}

// func GetFlagsmithClient() *flagsmith.Client {
// 	r := GetSingletonObject()
// 	return r.flagsmithClient
// }

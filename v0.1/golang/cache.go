package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type CacheDetails struct {
	Host        string        `json:"host"`
	IdleTimeout time.Duration `json:"idle_timeout"`
	MaxActive   int64         `json:"max_active"`
	MaxIdle     int64         `json:"max_idle"`
	Port        int64         `json:"port"`
	Protocol    string        `json:"protocol"`
}

type CacheInterface struct {
	pool *redis.Pool
}

const (
	DELIMITER = ":"
)

var (
	errNilConfig = errors.New(
		"redix.Create() - nil config provided",
	)
	errOneOrFewerArgs = errors.New("1 or fewer arguments provided")
)

func (ce *CacheInterface) Exec(args *[]interface{}, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	if len(*args) < 2 {
		return nil, errOneOrFewerArgs
	}

	conn := ce.pool.Get()
	defer conn.Close()

	upckdArgs := *args
	firstArg := fmt.Sprint(upckdArgs[0])
	restArgs := upckdArgs[1:]

	return conn.Do(firstArg, restArgs...)
}

func createRedisPool(config *CacheDetails) (*redis.Pool, error) {
	if config == nil {
		return nil, errNilConfig
	}

	redisAddress := fmt.Sprint(config.Host, DELIMITER, config.Port)

	pool := redis.Pool{
		MaxIdle:     int(config.MaxIdle),
		IdleTimeout: config.IdleTimeout,
		MaxActive:   int(config.MaxActive),
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				config.Protocol,
				redisAddress,
			)
		},
	}

	return &pool, nil
}

func NewInterface(details *CacheDetails) (cacheInterface *CacheInterface, err error) {
	pool, errPool := createRedisPool(details)
	if errPool != nil {
		return nil, errPool
	}

	cache := CacheInterface{
		pool: pool,
	}

	return &cache, nil
}

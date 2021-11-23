package cache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// now we have to create the redis connection
//
// probably a struct with an Exec()

type ConfigDetails struct {
	Filepath     string `json:"filepath"`
	FilepathTest string `json:"filepath_test"`
}

type CacheDetails struct {
	Host        string        `json:"host"`
	IdleTimeout time.Duration `json:"idle_timeout"`
	MaxActive   int64         `json:"max_active"`
	MaxIdle     int64         `json:"max_idle"`
	MaxSamples  int64         `json:"max_samples"`
	MaxSizeInMB string        `json:"max_size_in_mb"`
	Port        int64         `json:"port"`
	Protocol    string        `json:"protocol"`
}

type SuperCacheDetails struct {
	Config      ConfigDetails `json:"config"`
	Cache       CacheDetails  `json:"cache"`
}

const (
	DELIMITER = ":"
)

var (
	detailsPath         = os.Getenv("CONFIG_FILEPATH")
	Details, DetailsErr = ReadDetailsFromFile(detailsPath)

	errNilConfig = errors.New(
		"redix.Create() - nil config provided",
	)
	errOneOrFewerArgs = errors.New("1 or fewer arguments provided")

	pool, errPool = createRedisPool(&details.Details.Cache)
)

func readFile(path string) (*[]byte, error) {
	detailsJSON, errDetiailsJSON := ioutil.ReadFile(path)
	return &detailsJSON, errDetiailsJSON
}

func parseDetails(detailsJSON *[]byte, err error) (*SuperCacheDetails, error) {
	if err != nil {
		return nil, err
	}

	var details SuperCacheDetails
	errDetails := json.Unmarshal(*detailsJSON, &details)

	return &details, errDetails
}

func ReadDetailsFromFile(path string) (*SuperCacheDetails, error) {
	detailsJSON, errDetailsJSON := readFile(path)
	return parseDetails(detailsJSON, errDetailsJSON)
}

func createRedisPool(config *details.CacheDetails) (*redis.Pool, error) {
	if config == nil {
		return nil, errNilConfig
	}

	redisAddress := fmt.Sprint(config.Host, DELIMITER, config.Port)

	pool := redis.Pool{
		MaxIdle:     int(config.MaxIdle),
		IdleTimeout: config.IdleTimeout,
		MaxActive:   int(config.MaxActive),
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				config.Protocol,
				redisAddress,
			)

			if err != nil {
				return nil, err
			}

			return conn, nil
		},
	}

	return &pool, nil
}

func Exec(args *[]interface{}, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	if errPool != nil {
		return nil, errPool
	}

	if len(*args) < 2 {
		return nil, errOneOrFewerArgs
	}

	conn := pool.Get()
	defer conn.Close()

	upckdArgs := *args
	firstArg := fmt.Sprint(upckdArgs[0])
	restArgs := upckdArgs[1:]

	return conn.Do(firstArg, restArgs...)
}
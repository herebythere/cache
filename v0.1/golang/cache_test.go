package cache

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	testEntry       = "hello_world_test"
	testEntryResult = "how_are_you_starshine?"
)

// Env variables for tests
var (
	cacheAddress = os.Getenv("TEST_CACHE_HOST_ADDRESS")
	cachePort    = os.Getenv("TEST_CACHE_HOST_PORT")

	details = CacheDetails{
		Host:        "127.0.0.1",
		IdleTimeout: time.Second * 2,
		MaxActive:   128,
		MaxIdle:     3,
		Port:        3010,
		Protocol:    "tcp",
	}

	cache, errCache = NewInterface(&details)
)

func getDetails() (*CacheDetails, error) {
	if cacheAddress == "" {
		cacheAddress = "127.0.0.1"
	}
	if cachePort == "" {
		cachePort = "6379"
	}

	cachePortInt64, errCachePortInt64 := strconv.ParseInt(cachePort, 10, 64)
	if errCachePortInt64 != nil {
		return nil, errCachePortInt64
	}

	details := CacheDetails{
		Host:        cacheAddress,
		IdleTimeout: time.Second * 2,
		MaxActive:   128,
		MaxIdle:     3,
		Port:        cachePortInt64,
		Protocol:    "tcp",
	}

	return &details, nil
}

func TestExists(t *testing.T) {
	if cache == nil {
		t.Error("nil parameters should return nil")
	}
	if errCache != nil {
		t.Error("nil paramters should return error")
	}
}

func TestExec(t *testing.T) {
	setCommands := []interface{}{"SET", testEntry, testEntryResult}
	entry, errEntry := cache.Exec(&setCommands, nil)
	if errEntry != nil {
		t.Fail()
		t.Logf(errEntry.Error())
	}

	if entry == nil {
		t.Fail()
		t.Logf("setter.Set should retrun an entry")
	}

	getCommands := []interface{}{"GET", testEntry}
	getterEntry, errGetterEntry := redis.String(cache.Exec(&getCommands, nil))
	if errGetterEntry != nil {
		t.Fail()
		t.Logf(errGetterEntry.Error())
	}

	if getterEntry != testEntryResult {
		t.Fail()
		t.Logf("setter.Get should equal found count")
	}
}

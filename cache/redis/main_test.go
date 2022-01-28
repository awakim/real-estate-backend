package cache

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/awakim/immoblock-backend/config"
	"github.com/go-redis/redis/v8"
)

var testCache Cache

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testRDB := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
	})

	testCache = NewCache(testRDB)

	os.Exit(m.Run())
}

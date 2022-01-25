package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/awakim/immoblock-backend/api"
	cache "github.com/awakim/immoblock-backend/cache/redis"
	db "github.com/awakim/immoblock-backend/db/sqlc"
	"github.com/awakim/immoblock-backend/util"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("cannot connect to redis:", err)
	}

	store := db.NewStore(conn)
	cache := cache.NewCache(rdb)

	server, err := api.NewServer(config, store, cache)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/awakim/immoblock-backend/api"
	cache "github.com/awakim/immoblock-backend/cache/redis"
	"github.com/awakim/immoblock-backend/config"
	db "github.com/awakim/immoblock-backend/db/sqlc"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func main() {

	ts := time.Now()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config, err := config.LoadConfig(".")
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

	srv := &http.Server{
		Addr:         server.Config.ServerAddress,
		Handler:      server.Router,
		WriteTimeout: 3000 * time.Millisecond,
		ReadTimeout:  3000 * time.Millisecond,
		IdleTimeout:  3000 * time.Millisecond,
	}
	log.Printf("Time for startup: %v", time.Since(ts))
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	err = conn.Close()
	if err != nil {
		log.Fatalf("could not close DB connection %v", err)
	}

	err = rdb.Close()
	if err != nil {
		log.Fatalf("could not close Redis connection %v", err)
	}
	log.Println("Closed DB connection.")
	log.Println("Closed Redis connection.")
	log.Println("Shutting down gracefully. Send signal again to force kill.")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting.")
}

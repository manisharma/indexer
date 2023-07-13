package main

import (
	"context"
	"indexer/pkg/config"
	"indexer/pkg/db"
	"indexer/pkg/handler"
	"indexer/pkg/indexer"
	"indexer/pkg/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// initialize a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// parse app's config
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("config.Parse() failed, err: %v\n", err.Error())
	}

	// connect database
	pool, err := db.Connect(ctx, cfg.Postgres)
	if err != nil {
		log.Fatalf("db.Connect() failed, err: %v\n", err.Error())
	}
	defer pool.Close()

	// run migrations
	err = db.Migrate(cfg.Postgres)
	if err != nil {
		log.Fatalf("db.Migrate() failed, err: %v\n", err.Error())
	}

	// create data store
	repo := store.New(pool)

	// create new instance of indexer
	chain, err := indexer.New(ctx, cfg.ClientURL)
	if err != nil {
		log.Fatalf("indexer.New() failed, err: %v\n", err.Error())
	}

	// subscribe to incoming Epochs
	epochStream := chain.SubscribeToSlots(ctx)
	go func(ctx context.Context) {
		for epochResult := range epochStream {
			if epochResult.Error != nil {
				log.Printf("subscription failed, err: %v\n", epochResult.Error.Error())
			} else if epochResult.Epoch != nil {
				err = repo.Create(ctx, *epochResult.Epoch)
				if err != nil {
					log.Printf("repo.Create() failed, err: %v\n", err.Error())
				}
				err = repo.KeepOnlyTop5(ctx, epochResult.Epoch.EpochNumber)
				if err != nil {
					log.Printf("repo.KeepOnlyTop5() failed, err: %v\n", err.Error())
				}
			}
		}
	}(ctx)

	// initialize http handler
	httpHandler := handler.New(repo)

	// initialize http server
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      httpHandler,
	}
	// start server
	go func() {
		log.Println("server listening on http://localhost:8080")
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to listen, err: %v", err.Error())
		}
	}()

	// await termination/interruption for graceful shutdown
	interruption := make(chan os.Signal, 1)
	signal.Notify(interruption, syscall.SIGABRT, syscall.SIGTERM)
	<-interruption
	log.Println("interrupted, shutdown initiated")
	ctxWithTimeOut, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	server.Shutdown(ctxWithTimeOut)
	log.Println("graceful shutdown complete")
}

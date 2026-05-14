package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/omurilo/shareless/api/handler"
	"github.com/omurilo/shareless/internal/database"
	server "github.com/omurilo/shareless/internal/http"
	"github.com/omurilo/shareless/internal/telemetry"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file")
	}

	PORT, ok := os.LookupEnv("PORT")
	if !ok {
		PORT = "3000"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdown, err := telemetry.Setup(ctx)
	if err != nil {
		log.Printf("Failed to setup telemetry: %v", err)
	} else {
		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := shutdown(shutdownCtx); err != nil {
				log.Printf("Failed to shutdown telemetry: %v", err)
			}
		}()
	}

	db := database.NewDbClient()
	sh := handler.NewShareHandler(db)
	shr := handler.NewSharedHandler(db)

	httpServer := server.NewHttpServer(sh, shr)

	log.Printf("Executing server on port: %s", PORT)
	panic(http.ListenAndServe(fmt.Sprintf(":%s", PORT), httpServer.Instance))
}

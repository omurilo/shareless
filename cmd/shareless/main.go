package main

import (
	"context"
	"fmt"
	"log/slog"
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
		slog.Warn("Failed to load .env file")
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	PORT, ok := os.LookupEnv("PORT")
	if !ok {
		PORT = "3000"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdown, err := telemetry.Setup(ctx)
	if err != nil {
		slog.Error("Failed to setup telemetry", slog.Any("error", err))
	} else {
		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := shutdown(shutdownCtx); err != nil {
				slog.Error("Failed to shutdown telemetry", slog.Any("error", err))
			}
		}()
	}

	db := database.NewDbClient()
	sh := handler.NewShareHandler(db)
	shr := handler.NewSharedHandler(db)

	httpServer := server.NewHttpServer(sh, shr)

	slog.Info("Executing server", slog.String("port", PORT))
	panic(http.ListenAndServe(fmt.Sprintf(":%s", PORT), httpServer.Instance))
}

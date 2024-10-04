package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/omurilo/shareless/api/handler"
	"github.com/omurilo/shareless/internal/database"
	"github.com/omurilo/shareless/internal/http"
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

	db := database.NewDbClient()
	sh := handler.NewShareHandler(db)
	shr := handler.NewSharedHandler(db)

	httpServer := server.NewHttpServer(sh, shr)

	log.Printf("Executing server on port: %s", PORT)
	panic(http.ListenAndServe(fmt.Sprintf(":%s", PORT), httpServer.Instance))
}

package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Duration string

const (
	FiveMinutes Duration = "5m"
	OneHour     Duration = "1h"
	ThreeHours  Duration = "3h"
	SevenHours  Duration = "7h"
	OneDay      Duration = "24h"
)

type ShareInput struct {
	Duration       Duration `json:"duration"`
	ExpireOnOpened string   `json:"expire_on_opened"`
	PublicKey      string   `json:"publicKey"`
	PrivateKey     string   `json:"privateKey"`
}

type SharedDocument struct {
	Id             uuid.UUID `json:"id"`
	PublicKey      string    `json:"-"`
	PrivateKey     string    `json:"-"`
	ExpireOnOpened bool      `json:"-"`
	Host           string    `json:"host"`
}

type ShareHandler struct {
	db *redis.Client
}

func init() {
	_, ok := os.LookupEnv("HOST_URL")
	if !ok {
		log.Println("HOST_URL is not defined, consider setting to get a full shareable link!")
	}
}

func NewShareHandler(db *redis.Client) *ShareHandler {
	return &ShareHandler{db}
}

func (s *ShareHandler) Share(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body ShareInput
	err := d.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = validateDuration(body.Duration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	duration, err := time.ParseDuration(string(body.Duration))

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	shared := SharedDocument{
		Id:             uuid.New(),
		PublicKey:      body.PublicKey,
		PrivateKey:     body.PrivateKey,
		ExpireOnOpened: body.ExpireOnOpened == "on",
		Host:           os.Getenv("HOST_URL"),
	}

	err = s.db.HSet(r.Context(), shared.Id.String(), "publicKey", shared.PublicKey, "privateKey", shared.PrivateKey, "expire_on_opened", shared.ExpireOnOpened).Err()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "An error has ocurred on generate a shareless", http.StatusInternalServerError)
	}

	err = s.db.Expire(r.Context(), shared.Id.String(), duration).Err()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "An error has ocurred on generate a shareless", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": fmt.Sprintf("%s/shared/%s", shared.Host, shared.Id.String()),
	})
}

func validateDuration(duration Duration) error {
	switch duration {
	case FiveMinutes:
	case OneHour:
	case ThreeHours:
	case SevenHours:
	case OneDay:
		break
	default:
		return errors.New(
			fmt.Sprintf("Duration is invalid, only valid durations is: %v", []Duration{FiveMinutes, OneHour, ThreeHours, SevenHours, OneDay}),
		)
	}

	return nil
}

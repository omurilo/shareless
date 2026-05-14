package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/omurilo/shareless/pkg/cipher"
	"github.com/omurilo/shareless/web"
	"github.com/redis/go-redis/v9"
	"github.com/x-way/crawlerdetect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type SharedHandler struct {
	db           *redis.Client
	sharesViewed metric.Int64Counter
}

func NewSharedHandler(db *redis.Client) *SharedHandler {
	meter := otel.GetMeterProvider().Meter("github.com/omurilo/shareless")
	sharesViewed, _ := meter.Int64Counter(
		"shareless.shares.viewed",
		metric.WithDescription("Total number of shares viewed"),
		metric.WithUnit("{view}"),
	)
	return &SharedHandler{db, sharesViewed}
}

func (s *SharedHandler) Shared(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	token := r.URL.Query().Get("token")

	uastring := r.Header.Get("User-Agent")

	if crawlerdetect.IsCrawler(uastring) {
		w.Header().Set("Content-Type", "text/html")
		web.Shared(w, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	hash := s.db.HGet(r.Context(), id, "hash")
	if hash.Err() != nil {
		http.Error(w, "The document not found or its expired", http.StatusNotFound)
		return
	}

	expire := s.db.HGet(r.Context(), id, "expire_on_opened")
	if expire.Err() != nil {
		http.Error(w, "The document not found or its expired", http.StatusNotFound)
		return
	}

	if expire.Val() == "1" {
		expireCmd := s.db.Expire(r.Context(), id, 0)
		if expireCmd.Err() != nil {
			slog.ErrorContext(r.Context(), "failed to expire secret", slog.String("id", id))
		}
	}

	plainText, err := cipher.Decrypter(token, hash.Val())
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to decrypt share", slog.Any("error", err))
		http.Error(w, "The token was sent is invalid", http.StatusUnprocessableEntity)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html")
		web.Shared(w, map[string]string{"Text": plainText})
		s.sharesViewed.Add(r.Context(), 1, metric.WithAttributes(attribute.Bool("expire_on_opened", expire.Val() == "1")))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"text": plainText})
	s.sharesViewed.Add(r.Context(), 1, metric.WithAttributes(attribute.Bool("expire_on_opened", expire.Val() == "1")))
}

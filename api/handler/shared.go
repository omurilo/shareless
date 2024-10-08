package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/omurilo/shareless/web"
	"github.com/redis/go-redis/v9"
	"github.com/x-way/crawlerdetect"
)

type SharedHandler struct {
	db *redis.Client
}

func NewSharedHandler(db *redis.Client) *SharedHandler {
	return &SharedHandler{db}
}

func (s *SharedHandler) Shared(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	uastring := r.Header.Get("User-Agent")

	if crawlerdetect.IsCrawler(uastring) {
		w.Header().Set("Content-Type", "text/html")
		web.Shared(w, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	privateKey := s.db.HGet(r.Context(), id, "privateKey")
	if privateKey.Err() != nil {
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
			log.Printf("Error when expire secret: %v\n", id)
		}
	}

	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html")
		web.Shared(w, map[string]string{"PrivateKey": privateKey.Val()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"privateKey": privateKey.Val()})
}

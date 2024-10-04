package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/omurilo/shareless/pkg/cipher"
	"github.com/omurilo/shareless/web"
	"github.com/redis/go-redis/v9"
)

type SharedHandler struct {
	db *redis.Client
}

func NewSharedHandler(db *redis.Client) *SharedHandler {
	return &SharedHandler{db}
}

func (s *SharedHandler) Shared(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	token := r.URL.Query().Get("token")

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
			log.Printf("Error when expire secret: %v\n", id)
		}
	}

	plainText, err := cipher.Decrypter(token, hash.Val())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "The token was sent is invalid", http.StatusUnprocessableEntity)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html")
		web.Shared(w, map[string]string{"Text": plainText})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"text": plainText})
}

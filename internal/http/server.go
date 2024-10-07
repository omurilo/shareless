package server

import (
	"net/http"

	"github.com/omurilo/shareless/api/handler"
	"github.com/omurilo/shareless/web"
)

type HttpServer struct {
	Instance *http.ServeMux
}

func NewHttpServer(sh *handler.ShareHandler, shr *handler.SharedHandler) *HttpServer {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /shared/{id}", shr.Shared)
	mux.HandleFunc("POST /share", sh.Share)
	mux.HandleFunc("GET /", web.Shareless)
	mux.HandleFunc("GET /about", web.About)
	mux.HandleFunc("GET /privacy", web.Privacy)

	return &HttpServer{mux}
}

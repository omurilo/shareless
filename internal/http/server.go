package server

import (
	"net/http"

	"github.com/omurilo/shareless/api/handler"
	"github.com/omurilo/shareless/web"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HttpServer struct {
	Instance http.Handler
}

func NewHttpServer(sh *handler.ShareHandler, shr *handler.SharedHandler) *HttpServer {
	mux := http.NewServeMux()

	mux.Handle("GET /shared/{id}", otelhttp.NewHandler(http.HandlerFunc(shr.Shared), "shared.get"))
	mux.Handle("POST /share", otelhttp.NewHandler(http.HandlerFunc(sh.Share), "share.create"))
	mux.HandleFunc("GET /", web.Shareless)
	mux.HandleFunc("GET /about", web.About)
	mux.HandleFunc("GET /privacy", web.Privacy)

	return &HttpServer{mux}
}

package server

import (
	"net/http"

	"github.com/omurilo/shareless/api/handler"
	"github.com/omurilo/shareless/internal/middleware"
	"github.com/omurilo/shareless/web"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HttpServer struct {
	Instance http.Handler
}

func NewHttpServer(sh *handler.ShareHandler, shr *handler.SharedHandler) *HttpServer {
	mux := http.NewServeMux()

	apiHandler := func(h http.Handler, op string) http.Handler {
		return otelhttp.NewHandler(middleware.Logger(h), op)
	}

	mux.Handle("GET /shared/{id}", apiHandler(http.HandlerFunc(shr.Shared), "shared.get"))
	mux.Handle("POST /share", apiHandler(http.HandlerFunc(sh.Share), "share.create"))
	mux.HandleFunc("GET /", web.Shareless)
	mux.HandleFunc("GET /about", web.About)
	mux.HandleFunc("GET /privacy", web.Privacy)

	return &HttpServer{mux}
}

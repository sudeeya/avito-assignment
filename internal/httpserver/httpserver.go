package httpserver

import (
	"net/http"
	"strconv"

	"github.com/sudeeya/avito-assignment/internal/config"
)

func NewServer(cfg config.ServerConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.ServerHTTPPort),
		Handler: handler,
	}
}

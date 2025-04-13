package httpserver

import (
	"net/http"
	"strconv"

	"github.com/sudeeya/avito-assignment/internal/config"
	v1 "github.com/sudeeya/avito-assignment/internal/controller/http/v1"
	"github.com/sudeeya/avito-assignment/internal/service"
)

func NewServer(cfg config.ServerConfig, services *service.Services) *http.Server {
	return &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.ServerHTTPPort),
		Handler: v1.NewRouter(services),
	}
}

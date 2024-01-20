package http_controller

import (
	"chatmerger/internal/usecase"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type HttpController struct {
	sh *http.Server
	requiredUsecases
}

type Config struct {
	Host string
	Port int
}

type requiredUsecases interface {
	usecase.CreateClientUc
	usecase.DeleteClientUc
	usecase.ClientsUc
}

func NewHttpController(cfg Config, usecases requiredUsecases) *HttpController {
	var router = mux.NewRouter()

	httpServer := &http.Server{
		Addr:           cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var adminServer = &HttpController{
		sh:               httpServer,
		requiredUsecases: usecases,
	}
	adminServer.registerHttpServerRoutes(router)
	return adminServer
}

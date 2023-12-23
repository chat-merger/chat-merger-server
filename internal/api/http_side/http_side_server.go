package http_side

import (
	"chatmerger/internal/usecase"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type HttpSideServer struct {
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
	usecase.ClientsListUc
	usecase.ConnectedClientsListUc
}

func NewHttpSideServer(cfg Config, usecases requiredUsecases) *HttpSideServer {
	var router = mux.NewRouter()

	httpServer := &http.Server{
		Addr:           cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var adminServer = &HttpSideServer{
		sh:               httpServer,
		requiredUsecases: usecases,
	}
	adminServer.registerHttpServerRoutes(router)
	return adminServer
}

func (s *HttpSideServer) Serve(ctx context.Context) error {
	go s.contextCancelHandler(ctx)

	if err := s.sh.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("ListenAndServe http server: %s", err)
	}
	return nil
}

func (s *HttpSideServer) contextCancelHandler(ctx context.Context) {
	select {
	case <-ctx.Done():
		s.sh.Shutdown(context.Background())
	}
}

func (s *HttpSideServer) registerHttpServerRoutes(router *mux.Router) {

	router.HandleFunc("/", s.index)

	//var apiRoutes = router.PathPrefix("/api")
	router.Path("/api").HandlerFunc(s.createClientHandler).Methods(http.MethodPost)
	router.Path("/api/{id}").HandlerFunc(s.deleteClientHandler).Methods(http.MethodDelete)
	router.Path("/api").HandlerFunc(s.getClientsHandler).Methods(http.MethodGet)
}

package admin

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

type Server struct {
	sh *http.Server
	Usecases
}

type Config struct {
	Host string
	Port int
}

type Usecases struct {
	usecase.CreateClientUc
	usecase.DeleteClientUc
	usecase.ClientsListUc
	usecase.ConnectedClientsListUc
}

func NewAdminPanelServer(cfg Config, usecases Usecases) *Server {
	var router = mux.NewRouter()

	httpServer := &http.Server{
		Addr:           cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var adminServer = &Server{
		sh:       httpServer,
		Usecases: usecases,
	}
	adminServer.registerHttpServerRoutes(router)
	return adminServer
}

func (s *Server) Serve(ctx context.Context) error {
	go s.contextCancelHandler(ctx)

	if err := s.sh.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("ListenAndServe http server: %s", err)
	}
	return nil
}

func (s *Server) contextCancelHandler(ctx context.Context) {
	select {
	case <-ctx.Done():
		s.sh.Shutdown(context.Background())
	}
}

func (s *Server) registerHttpServerRoutes(router *mux.Router) {

	router.HandleFunc("/", s.index)

	//var apiRoutes = router.PathPrefix("/api")
	router.Path("/api").HandlerFunc(s.createClientHandler).Methods(http.MethodPost)
	router.Path("/api/{id}").HandlerFunc(s.deleteClientHandler).Methods(http.MethodDelete)
	router.Path("/api").HandlerFunc(s.getClientsHandler).Methods(http.MethodGet)
}

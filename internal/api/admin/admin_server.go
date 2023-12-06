package admin

import (
	"chatmerger/internal/usecase"
	"context"
	"errors"
	"fmt"
	"net/http"
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
	usecase.ClientsSessionsListUc
}

func NewAdminServer(cfg Config, usecases Usecases) *Server {
	var mux = http.NewServeMux()

	httpServer := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var adminServer = &Server{
		sh:       httpServer,
		Usecases: usecases,
	}
	adminServer.registerHttpServerRoutes(mux)
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

func (s *Server) registerHttpServerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		var b []byte
		request.Body.Read(b)
		writer.Write(b)
	})
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			s.createClientHandler(writer, request)
		} else if request.Method == http.MethodDelete {
			s.deleteClientHandler(writer, request)
		} else if request.Method == http.MethodGet {
			s.getClientsHandler(writer, request)
		} else {
			writer.WriteHeader(http.StatusNotFound)
		}
	})
}

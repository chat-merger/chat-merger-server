package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	//uc uc.Usecases
	sh *http.Server
}

type Config struct {
	Host string
	Port int
}

func NewAdminServer(cfg Config) *Server {
	var mux = http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello"))
	})
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		var b []byte
		request.Body.Read(b)
		writer.Write(b)
	})
	httpServer := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	var adminServer = &Server{
		//uc: usecases,
		sh: httpServer,
	}
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

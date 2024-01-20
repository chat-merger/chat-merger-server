package http_controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *HttpController) Run(ctx context.Context) error {
	go s.contextCancelHandler(ctx)

	if err := s.sh.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("ListenAndServe http server: %s", err)
	}
	return nil
}

func (s *HttpController) contextCancelHandler(ctx context.Context) {
	select {
	case <-ctx.Done():
		s.sh.Shutdown(context.Background())
	}
}

func (s *HttpController) registerHttpServerRoutes(router *mux.Router) {

	router.HandleFunc("/", s.index)

	//var apiRoutes = router.PathPrefix("/api")
	router.Path("/api").HandlerFunc(s.createClientHandler).Methods(http.MethodPost)
	router.Path("/api/{id}").HandlerFunc(s.deleteClientHandler).Methods(http.MethodDelete)
	router.Path("/api").HandlerFunc(s.getClientsHandler).Methods(http.MethodGet)
}

package app

import (
	adm "chatmerger/internal/api/admin"
	"chatmerger/internal/api/pb"
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/usecase"
	. "chatmerger/internal/repositories/client_sessions_repository"
	. "chatmerger/internal/repositories/clients_repository"
	"context"
	"errors"
	"log"
	"net/http"
)

type State uint8

const (
	Working State = iota
	Running
	Stopping
	Stopped
	RunningFailure
	StoppingFailure
	InternalError
)

type application struct {
	clientRepository  domain.ClientRepository
	sessionRepository domain.ClientsSessionRepository
	apiServer         domain.Handler
	adminServer       domain.Handler
}

func (a *application) stop() {
	// todo: a.changeStatus(Stopping)
	cc, err := a.sessionRepository.Connected()
	if err == nil {
		for _, client := range cc {
			a.sessionRepository.Disconnect(client.Id)
		}
	}
}

func Run(ctx context.Context) error {
	var sessionsRepo = &ClientSessionsRepositoryBase{}
	var clientsRepo = &ClientsRepositoryBase{}
	var apiServer = pb.NewApiServer(pb.Config{
		Host: "localhost",
		Port: 8080,
	})
	var usecases = usecase.Usecases{}
	var adminServer = adm.NewAdminServer(usecases, adm.Config{
		Host: "localhost",
		Port: 8081,
	})

	var app = &application{
		clientRepository:  clientsRepo,
		sessionRepository: sessionsRepo,
		apiServer:         apiServer,
		adminServer:       adminServer,
	}

	go app.startAdminServer(ctx)
	go app.startApiServer(ctx)

	return nil
}

func (a *application) contextCancelHandler(ctx context.Context) {
	select {
	case <-ctx.Done():
		a.stop()
	}
}

func (a *application) startAdminServer(ctx context.Context) {
	err := a.adminServer.Serve(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("admin server serve: %s", err)
	}
}

func (a *application) startApiServer(ctx context.Context) {
	err := a.apiServer.Serve(ctx)
	if err != nil {
		log.Fatalf("api server serve: %s", err)
	}
}

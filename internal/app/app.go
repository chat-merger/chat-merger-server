package app

import (
	"chatmerger/internal/api"
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/usecase"
	. "chatmerger/internal/repositories/client_sessions_repository"
	. "chatmerger/internal/repositories/clients_repository"
	"context"
	"fmt"
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
	apiServer         domain.ApiServer
}

func (a *application) stop() {
	// todo: a.changeStatus(Stopping)
	cc, err := a.sessionRepository.Connected()
	if err == nil {
		for _, client := range cc {
			a.sessionRepository.Disconnect(client.Id)
		}
	}
	a.apiServer.Stop()
}

func Run(ctx context.Context) error {
	var sessionsRepo = &ClientSessionsRepositoryBase{}
	var clientsRepo = &ClientsRepositoryBase{}
	var usecases = usecase.Usecases{}
	var apiServer = api.NewApiServer(usecases, api.Config{
		Host: "localhost",
		Port: 8080,
	})

	var app = &application{
		clientRepository:  clientsRepo,
		sessionRepository: sessionsRepo,
		apiServer:         apiServer,
	}

	go func() {
		select {
		case <-ctx.Done():
			app.stop()
		}
	}()

	err := app.apiServer.Start()
	if err != nil {
		return fmt.Errorf("start api server failed: %s", err)
	}

	return nil
}

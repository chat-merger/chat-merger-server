package app

import (
	adm "chatmerger/internal/api/admin"
	"chatmerger/internal/api/pb"
	"chatmerger/internal/domain"
	. "chatmerger/internal/repositories/client_sessions_repository"
	. "chatmerger/internal/repositories/clients_repository"
	"chatmerger/internal/uc"
	"chatmerger/internal/usecase"
	"context"
	"log"
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
	apiHandler        domain.Handler
	adminHandler      domain.Handler
	usecases          *Usecases
}

type Usecases struct {
	usecase.SendMessageToEveryoneExceptUc
	usecase.CreateClientSessionUc
	usecase.DropClientSessionUc
	usecase.ClientsListUc
	usecase.ClientsSessionsListUc
	usecase.CreateClientUc
	usecase.DeleteClientUc
}

// use for graceful shutdown
func (a *application) shutdown() {
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
	var usecases = &Usecases{
		SendMessageToEveryoneExceptUc: &uc.SendMessageToEveryoneExcept{},
		CreateClientSessionUc:         &uc.CreateClientSession{},
		DropClientSessionUc:           &uc.DropClientSession{},
		ClientsListUc:                 &uc.ClientsList{},
		ClientsSessionsListUc:         &uc.ClientsSessionsList{},
		CreateClientUc:                &uc.CreateClient{},
		DeleteClientUc:                &uc.DeleteClient{},
	}

	// create and run clients api handler
	var apiHandler = pb.NewApiServer(
		pb.Config{
			Host: "localhost",
			Port: 8080,
		},
		pb.Usecases{
			SendMessageToEveryoneExceptUc: usecases,
			CreateClientSessionUc:         usecases,
			DropClientSessionUc:           usecases,
		},
	)
	go serveHandler(apiHandler, ctx)

	// crate and run admin api handler
	var adminHandler = adm.NewAdminServer(usecases, adm.Config{
		Host: "localhost",
		Port: 8081,
	})
	go serveHandler(adminHandler, ctx)

	// todo:
	var _ = &application{
		clientRepository:  clientsRepo,
		sessionRepository: sessionsRepo,
		apiHandler:        apiHandler,
		adminHandler:      adminHandler,
	}

	return nil
}

func (a *application) contextCancelHandler(ctx context.Context) {
	select {
	case <-ctx.Done():
		a.shutdown()
	}
}

func serveHandler(h domain.Handler, ctx context.Context) {
	err := h.Serve(ctx)
	if err != nil {
		log.Fatalf("handler serve: %s", err)
	}
}

package app

import (
	"chatmerger/internal/domain"
	"context"
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
	apiServer          domain.ApiServer
}

func (a *application) stop() {
	// todo: a.changeStatus(Stopping)
	a.apiServer.Stop()
	cc, err := a.sessionRepository.Connected()
	if err == nil {
		for _, client := range cc {
			a.sessionRepository.Disconnect(client.Id)
		}
	}
}

func Run(ctx context.Context) error {
	var adapterRepo =

	var app = &application{
		clientRepository:  nil,
		sessionRepository: nil,
		apiServer:         nil,
	}

	go func() {
		select {
		case <-ctx.Done():
			app.stop()
		}
	}()

	return nil
}

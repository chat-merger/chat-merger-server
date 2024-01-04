package uc

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/usecase"
)

var _ usecase.ConnectedClientsListUc = (*ConnectedClientsList)(nil)

type ConnectedClientsList struct {
	sessionsRepo domain.ClientSessionsRepository
}

func NewConnectedClientsList(sessionsRepo domain.ClientSessionsRepository) *ConnectedClientsList {
	return &ConnectedClientsList{sessionsRepo: sessionsRepo}
}

func (r *ConnectedClientsList) ConnectedClientsList() ([]model.Client, error) {
	return r.sessionsRepo.Connected()
}

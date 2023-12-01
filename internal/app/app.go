package app

import (
	"chatmerger/internal/api"
	"chatmerger/internal/config"
	"chatmerger/internal/domain"
)

type App struct {
	cfg    domain.Config
	ccRepo domain.ClientConnectionRepository
}

func NewApp() *App {
	return &App{
		cfg:    &config.Config{},
		ccRepo: &api.ClientConnectRepositoryBase{},
	}
}

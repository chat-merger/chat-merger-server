package app

import (
	"chatmerger/internal/api/grpc_side"
	"chatmerger/internal/api/http_side"
	"chatmerger/internal/domain"
	"chatmerger/internal/usecase"
	"context"
)

type application struct {
	commonDeps
	httpSideCfg http_side.Config
	grpcSideCfg grpc_side.Config
}

type commonDeps struct {
	usecases *usecasesImpls
	ctx      context.Context
}

type usecasesImpls struct {
	usecase.CreateAndSendMsgToEveryoneExceptUc
	usecase.CreateClientSessionUc
	usecase.DropClientSessionUc
	usecase.ClientsListUc
	usecase.ConnectedClientsListUc
	usecase.CreateClientUc
	usecase.DeleteClientUc
}

type repositories struct {
	clientsRepo  domain.ClientsRepository
	sessionsRepo domain.ClientSessionsRepository
}

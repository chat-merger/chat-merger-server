package app

import (
	"chatmerger/internal/api/grpc_side"
	"chatmerger/internal/api/http_side"
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/config"
	csr "chatmerger/internal/repositories/client_sessions_repository"
	cr "chatmerger/internal/repositories/clients_repository"
	"chatmerger/internal/uc"
	"context"
	"fmt"
	"log"
)

func Run(ctx context.Context, cfg *config.Config) error {
	repos, err := initRepositories(cfg)
	if err != nil {
		return fmt.Errorf("init repositories: %s", err)
	}
	log.Println(msgs.RepositoriesInitialized)

	var usecases = newUsecases(repos)
	log.Println(msgs.UsecasesCreated)

	deps := commonDeps{usecases: usecases, ctx: ctx}
	app, errCh := newApplication(deps, cfg)

	// create and run clients api handler
	go app.runGrpcSideServer()
	// crate and run admin panel api handler
	go app.runHttpSideServer()

	log.Println(msgs.ApplicationStarted)

	return app.gracefulShutdownApplication(errCh)
}

func (a *application) gracefulShutdownApplication(errCh <-chan error) error {
	var err error
	select {
	case <-a.ctx.Done():
		log.Println(msgs.ApplicationReceiveCtxDone)
	case err = <-errCh:
		a.cancelFunc()
		log.Println(msgs.ApplicationReceiveInternalError)
	}
	a.wg.Wait()
	return err
}

func (a *application) runHttpSideServer() {
	a.wg.Add(1)
	defer a.wg.Done()
	log.Println(msgs.RunHttpSideServer)
	h := http_side.NewHttpSideServer(a.httpSideCfg, a.usecases)
	err := h.Serve(a.ctx)
	if err != nil {
		a.errorf("http side server serve: %s", err)
	}
	log.Println(msgs.StoppedHttpSideServer)
}

func (a *application) runGrpcSideServer() {
	a.wg.Add(1)
	defer a.wg.Done()
	log.Println(msgs.RunGrpcSideServer)
	h := grpc_side.NewGrpcSideServer(a.grpcSideCfg, a.usecases)
	err := h.Serve(a.ctx)
	if err != nil {
		a.errorf("grpc side server serve: %s", err)
	}
	log.Println(msgs.StoppedGrpcSideServer)
}

func (a *application) errorf(format string, args ...any) {
	select {
	case a.errCh <- fmt.Errorf(format, args):
	default:
	}
}

func initRepositories(cfg *config.Config) (*repositories, error) {
	sessionsRepo := csr.NewClientSessionsRepositoryBase()
	clientsRepo, err := cr.NewClientsRepositoryBase(cr.Config{
		FilePath: cfg.ClientsCfgFile,
	})
	if err != nil {
		return nil, fmt.Errorf("create clients repository: %s", err)
	}
	return &repositories{
		clientsRepo:  clientsRepo,
		sessionsRepo: sessionsRepo,
	}, nil

}

func newUsecases(repos *repositories) *usecasesImpls {
	return &usecasesImpls{
		ConnectedClientsListUc: uc.NewConnectedClientsList(repos.sessionsRepo),
		// clients api server
		CreateAndSendMsgToEveryoneExceptUc: uc.NewCreateAndSendMsgToEveryoneExcept(repos.sessionsRepo),
		CreateClientSessionUc:              uc.NewCreateClientSession(repos.clientsRepo, repos.sessionsRepo),
		DropClientSessionUc:                uc.NewDropClientSession(repos.sessionsRepo),
		// admin panel api server
		ClientsListUc:  uc.NewClientsList(repos.clientsRepo),
		CreateClientUc: uc.NewCreateClient(repos.clientsRepo),
		DeleteClientUc: uc.NewDeleteClient(repos.clientsRepo),
	}
}

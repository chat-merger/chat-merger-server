package app

import (
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/config"
	"chatmerger/internal/data/api/grpc_side"
	"chatmerger/internal/data/api/http_side"
	csr "chatmerger/internal/data/repositories/client_sessions_repository"
	cr "chatmerger/internal/data/repositories/clients_repository"
	"chatmerger/internal/data/uc"
	"chatmerger/internal/domain"
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

	deps := commonDeps{
		usecases: usecases,
		ctx:      ctx,
	}
	app, errCh := newApplication(ctx, deps)
	// create and run clients api handler
	gc := grpc_side.NewGrpcController(grpc_side.Config{
		Port: cfg.GrpcServerPort,
	}, app.usecases)
	go app.runController(gc, "GrpcController")

	// crate and run admin panel api handler
	hc := http_side.NewHttpController(http_side.Config{
		Port: cfg.HttpServerPort,
	}, app.usecases)
	go app.runController(hc, "HttpController")

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

func (a *application) runController(c domain.Controller, name string) {
	a.wg.Add(1)
	defer a.wg.Done()
	log.Println(msgs.RunController, name)
	err := c.Run(a.ctx)
	if err != nil {
		a.errorf("%s controller run: %s", name, err)
	}
	log.Println(msgs.StoppedController, name)
}

func (a *application) errorf(format string, args ...any) {
	select {
	case a.errCh <- fmt.Errorf(format, args...):
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

package app

import (
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/config"
	"chatmerger/internal/data/api/grpc_controller"
	"chatmerger/internal/data/api/http_controller"
	cr "chatmerger/internal/data/repositories/clients_repository"
	"chatmerger/internal/data/uc"
	"chatmerger/internal/domain"
	"chatmerger/internal/service/msgbus"
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

	ss := newServices()
	log.Println(msgs.ServicesCreated)

	var usecases = newUsecases(repos, ss)
	log.Println(msgs.UsecasesCreated)

	app, errCh := newApplication(ctx, usecases)
	// create and run clients api handler
	gc := grpc_controller.NewGrpcController(grpc_controller.Config{
		Port: cfg.GrpcServerPort,
	}, app.usecases)
	go app.runController(gc, "GrpcController")

	// crate and run admin panel api handler
	hc := http_controller.NewHttpController(http_controller.Config{
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
	clientsRepo, err := cr.NewClientsRepositoryBase(cr.Config{
		FilePath: cfg.ClientsCfgFile,
	})
	if err != nil {
		return nil, fmt.Errorf("create clients repository: %s", err)
	}
	return &repositories{
		cRepo: clientsRepo,
	}, nil

}

func newUsecases(r *repositories, ss *services) usecasesImpls {
	return usecasesImpls{
		// clients api server
		CreateAndSendMsgToEveryoneExceptUc: uc.NewCreateAndSendMsgToEveryoneExcept(r.cRepo, ss.bus),
		SubscribeClientToNewMsgsUc:         uc.NewSubscribeClientToNewMsgs(r.cRepo, ss.bus),
		DropClientSubscriptionUc:           uc.NewDropClientSubscription(r.cRepo, ss.bus),
		// admin panel api server
		ClientsUc:      uc.NewClients(r.cRepo),
		CreateClientUc: uc.NewCreateClient(r.cRepo),
		DeleteClientUc: uc.NewDeleteClient(r.cRepo),
	}
}

func newServices() *services {
	return &services{bus: msgbus.NewMessagesBus()}
}

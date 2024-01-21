package app

import (
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/component/eventbus"
	"chatmerger/internal/data/api/grpc_controller"
	"chatmerger/internal/data/api/http_controller"
	cr "chatmerger/internal/data/repositories/clients_repository"
	"chatmerger/internal/data/uc"
	"context"
	"fmt"
	"log"
)

func Run(ctx context.Context, cfg *Config) error {
	// init repos:
	repos, err := initRepositories(cfg)
	if err != nil {
		return fmt.Errorf("init repositories: %s", err)
	}
	log.Println(msgs.RepositoriesInitialized)

	// components:
	bus := eventbus.NewEventBus()
	log.Println(msgs.EventBusCreated)

	// create ucs:
	var usecases = newUsecases(repos, bus)
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

type Controller interface {
	Run(ctx context.Context) error
}

func (a *application) runController(c Controller, name string) {
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

func initRepositories(cfg *Config) (*repositories, error) {
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

func newUsecases(r *repositories, bus *eventbus.EventBus) usecasesImpls {
	return usecasesImpls{
		// clients api server
		CreateAndSendMsgToEveryoneExceptUc: uc.NewCreateAndSendMsgToEveryoneExcept(r.cRepo, bus),
		SubscribeClientToEventsUc:          uc.NewSubscribeClientToEvents(r.cRepo, bus),
		DropClientSubscriptionUc:           uc.NewDropClientSubscription(r.cRepo, bus),
		// admin panel api server
		ClientsUc:      uc.NewClients(r.cRepo),
		CreateClientUc: uc.NewCreateClient(r.cRepo),
		DeleteClientUc: uc.NewDeleteClient(r.cRepo),
	}
}

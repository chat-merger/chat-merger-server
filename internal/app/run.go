package app

import (
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/component/eventbus"
	"chatmerger/internal/component/sqlite"
	"chatmerger/internal/data/api/grpc_controller"
	"chatmerger/internal/data/api/http_controller"
	"chatmerger/internal/data/repositories/sqlite_clients_repo"
	"chatmerger/internal/data/uc"
	"context"
	"database/sql"
	"fmt"
	"log"
)

func Run(ctx context.Context, cfg *Config) error {
	// init db:
	db, err := sqlite.InitSqlite(sqlite.Config{
		DataSourceName: cfg.DbFile,
	})
	if err != nil {
		return fmt.Errorf("init databse: %s", err)
	}
	log.Println(msgs.DatabaseInitialized)

	// init repos:
	repos, err := initRepositories(db)
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

func initRepositories(db *sql.DB) (*repositories, error) {
	return &repositories{
		cRepo: sqlite_clients_repo.NewClientsRepository(db),
	}, nil

}

func newUsecases(r *repositories, bus *eventbus.EventBus) usecasesImpls {
	return usecasesImpls{
		// clients api server
		CreateAndSendMsgToEveryoneExceptUc: uc.NewCreateAndSendMsgToEveryoneExcept(r.cRepo, bus),
		SubscribeClientToEventsUc:          uc.NewSubscribeClientToEvents(r.cRepo, bus),
		DropClientSubscriptionUc:           uc.NewDropClientSubscription(r.cRepo, bus),
		// admin panel api server
		ClientsUc:      uc.NewClients(r.cRepo, bus),
		CreateClientUc: uc.NewCreateClient(r.cRepo),
		DeleteClientUc: uc.NewDeleteClient(r.cRepo, bus),
	}
}

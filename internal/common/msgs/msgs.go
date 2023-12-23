package msgs

import (
	"chatmerger/internal/common/vals"
	"fmt"
)

const (
	// main
	ServerStarting    = "Server Starting"
	ConfigInitialized = "Config Initialized"

	// application
	ApplicationStart          = "Start Application"
	ApplicationStarted        = "Application start is over, waiting when ctx done"
	UsecasesCreated           = "Usecases Created"
	RepositoriesCreated       = "Repositories Created"
	ApplicationReceiveCtxDone = "Application receive ctx.Done signal"

	// Grpc server
	RunGrpcSideServer       = "Run Grpc Side Server"
	ClientConnectedToServer = "Client Connected To Server"
	ClientSessionCreated    = "Client Session Created"

	// http server
	RunHttpSideServer = "Run Http Side Server"
)

var (
	ProgramWillForceExit = fmt.Sprintf("after %v seconds, the program will force exit\n", vals.GracefulShutdownTimeout.Seconds())
)

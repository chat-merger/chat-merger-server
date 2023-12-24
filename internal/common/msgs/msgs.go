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
	ApplicationStart                = "Start Application"
	ApplicationStarted              = "Application start is over, waiting when ctx done"
	UsecasesCreated                 = "Usecases Created"
	RepositoriesInitialized         = "RepositoriesInitialized"
	ApplicationReceiveCtxDone       = "Application receive ctx.Done signal"
	ApplicationReceiveInternalError = "ApplicationReceiveInternalError"

	// Grpc server
	RunGrpcSideServer            = "Run Grpc Side Server"
	StoppedGrpcSideServer        = "StoppedGrpcSideServer"
	ClientConnectedToServer      = "Client Connected To Server"
	ClientSessionCreated         = "Client Session Created"
	ClientSessionCloseConnection = "ClientSessionCloseConnection"
	NewMessageFromClient         = "NewMessageFromClient"

	// http server
	RunHttpSideServer     = "Run Http Side Server"
	StoppedHttpSideServer = "StoppedHttpSideServer"
)

var (
	ProgramWillForceExit = fmt.Sprintf("after %v seconds, the program will force exit", vals.GracefulShutdownTimeout.Seconds())
)

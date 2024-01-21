package msgs

const (
	// main
	ServerStarting    = "Server Starting"
	ConfigInitialized = "Config Initialized"

	// application
	ApplicationStart                = "Start Application"
	ApplicationStarted              = "Application start is over, waiting when ctx done"
	EventBusCreated                 = "EventBusCreated"
	UsecasesCreated                 = "UsecasesCreated"
	RepositoriesInitialized         = "RepositoriesInitialized"
	ApplicationReceiveCtxDone       = "Application receive ctx.Done signal"
	ApplicationReceiveInternalError = "ApplicationReceiveInternalError"

	// Grpc server
	RunGrpcController         = "RunGrpcController"
	StoppedGrpcController     = "StoppedGrpcController"
	ClientConnectedToServer   = "Client Connected To Server"
	ClientSubscribedToNewMsgs = "ClientSubscribedToNewMsgs"
	NewMessageFromClient      = "NewMessageFromClient"

	// http server
	RunHttpController     = "RunHttpController"
	StoppedHttpController = "StoppedHttpController"

	//  controller
	RunController     = "Run Controller"
	StoppedController = "Stopped Controller"
)

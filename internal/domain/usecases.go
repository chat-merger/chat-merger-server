package domain

type Usecases interface {
	// +onReciveMessageFromClient
	DropAllClientConnections()
	ClientsConnections()
	DropClientConnection(id int)
	ConnectClient()
	ClientsList()
	DeleteClient()
}

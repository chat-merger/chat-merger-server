package api


 type ClientRepositoryBase struct {}

	register(client Client) error
	registered() ([]Client, error)
	unregister(id int) error

	send(id int, msg Message) error
	connect(id int) (<-chan Message, error)
	disconnect(id int) error

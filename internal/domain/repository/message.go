package repository

type ClientRepository interface {
	register(client Client) error
	registered() ([]Client, error)
	unregister(id int) error

	// send(id int, msg Message) error
	// connect(id int, ch chan<- Message)
}

type Client struct {
	id     int
	name   string
	apiKey string
}

type MessengerRepository interface {
	send
}

type MailBox interface {
}

type Message struct {
}

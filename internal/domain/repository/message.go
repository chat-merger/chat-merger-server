package repository

type ClientRepository interface {
	register(client Client) error
	registered() ([]Client, error)
	unregister(id int) error

	send(id int, msg Message) error
	connect(id int) (<-chan Message, error)
	disconnect(id int) error
}

type Client struct {
	id     int
	name   string
	apiKey string
}

type ID string

type Message struct {
	id      ID
	replyId *ID
	date    int64
	author  *string
	from    string
	silent  bool
	body    Body
}

type Body interface {
	isBody()
}

type BodyText struct {
	format TextFormat
	value  string
}

type TextFormat string

const (
	Plain    TextFormat = "Plain"
	Markdown TextFormat = "Markdown"
)

type BodyMedia struct {
	kind    MediaType
	caption *string
	spoiler bool
	url     string
}

type MediaType string

const (
	Audio   MediaType = "Audio"
	Video   MediaType = "Video"
	File    MediaType = "File"
	Photo   MediaType = "Photo"
	Sticker MediaType = "Sticker"
)

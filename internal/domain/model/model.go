package model

type ClientSession struct {
	Client
	MsgCh <-chan Message
}

type Client struct {
	Id     int
	Name   string
	ApiKey string
}

type ID string

type Message struct {
	Id      ID
	ReplyId *ID
	Date    int64
	Author  *string
	From    string // adapter name
	Silent  bool
	Body    Body
}

type Body interface {
	IsBody()
}

type BodyText struct {
	Format TextFormat
	Value  string
}

type TextFormat string

const (
	Plain    TextFormat = "Plain"
	Markdown TextFormat = "Markdown"
)

type BodyMedia struct {
	Kind    MediaType
	Caption *string
	Spoiler bool
	Url     string
}

type MediaType string

const (
	Audio   MediaType = "Audio"
	Video   MediaType = "Video"
	File    MediaType = "File"
	Photo   MediaType = "Photo"
	Sticker MediaType = "Sticker"
)

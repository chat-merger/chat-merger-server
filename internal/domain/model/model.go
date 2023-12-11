package model

import (
	"fmt"
	"time"
)

// transfer data (dto) ( handler -> usecase(dto) -> dto to domain -> repository.meth(domain) = result)

type CreateClientSession struct {
	ApiKey ApiKey
}

type CreateClient struct {
	Name string `json:"name"`
}

// value models

type ID struct {
	value string
}

func (r ID) String() string {
	return fmt.Sprintf("%s", r.value)
}

func (r ID) Value() string {
	return r.value
}

func NewID(val string) ID {
	return ID{val}
}

type ApiKey struct {
	value string
}

func (r ApiKey) String() string {
	return fmt.Sprintf("%s", r.value)
}

func (r ApiKey) Value() string {
	return r.value
}

func NewApiKey(val string) ApiKey {
	return ApiKey{val}
}

// main models

type ClientSession struct {
	Client
	MsgCh <-chan Message
}

type Client struct {
	Id     ID
	Name   string
	ApiKey ApiKey
}

type Message struct {
	Id      ID
	ReplyId *ID
	Date    time.Time
	Author  *string
	From    string // client name
	Silent  bool
	Body    Body
}

type Body interface{ IsBody() }

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

package model

import (
	"fmt"
	"time"
)

// transfer data (dto) ( handler -> usecase(dto) -> dto to domain -> repository.meth(domain) = result)

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

type Client struct {
	Id     ID     `json:"id"`
	Name   string `json:"name,omitempty"`
	ApiKey ApiKey `json:"api_key"`
	Status ConnStatus
}

type ConnStatus uint8

const (
	_ ConnStatus = iota
	ConnStatusActive
	ConnStatusInactive
)

type ClientsFilter struct {
	Id     *ID
	Name   *string
	ApiKey *ApiKey
	Status *ConnStatus
}

type Message struct {
	Id       ID
	ReplyId  *ID
	Date     time.Time
	Username *string
	From     string // client name
	Silent   bool
	Body     Body
}

type Body interface{ IsBody() }

type BodyText struct {
	Format TextFormat
	Value  string
}

func (b *BodyText) IsBody() {}

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

func (b *BodyMedia) IsBody() {}

type MediaType string

const (
	Audio   MediaType = "Audio"
	Video   MediaType = "Video"
	File    MediaType = "File"
	Photo   MediaType = "Photo"
	Sticker MediaType = "Sticker"
)

// create message

type CreateMessage struct {
	ReplyId  *ID
	Date     time.Time
	Username *string
	From     string // client name
	Silent   bool
	Body     Body
}

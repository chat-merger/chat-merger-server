package model

import (
	"time"
)

// transfer data (dto) ( handler -> usecase(dto) -> dto to domain -> repository.meth(domain) = result)

type CreateClient struct {
	Name string `json:"name"`
}

// value models

type ApiKey string
type ID string

// main models

type ClientWithStatus struct {
	Id     ID
	Name   string
	ApiKey ApiKey
	Status ConnStatus
}

type Client struct {
	Id     ID
	Name   string
	ApiKey ApiKey
}

type ConnStatus uint8

const (
	ConnStatusUndefined ConnStatus = iota
	ConnStatusInactive
	ConnStatusActive
)

type ClientsFilter struct {
	Id     *ID
	Name   *string
	ApiKey *ApiKey
	Status ConnStatus
}

func (f ClientsFilter) ExceptStatus() ClientsFilterExceptStatus {
	return ClientsFilterExceptStatus{
		Id:     f.Id,
		Name:   f.Name,
		ApiKey: f.ApiKey,
	}
}

type ClientsFilterExceptStatus struct {
	Id     *ID
	Name   *string
	ApiKey *ApiKey
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

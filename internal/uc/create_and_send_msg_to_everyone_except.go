package uc

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/usecase"
	"fmt"
	"github.com/google/uuid"
	"slices"
)

var _ usecase.CreateAndSendMsgToEveryoneExceptUc = (*CreateAndSendMsgToEveryoneExcept)(nil)

type CreateAndSendMsgToEveryoneExcept struct {
	sessionsRepo domain.ClientsSessionRepository
}

func NewCreateAndSendMsgToEveryoneExcept(sessionsRepo domain.ClientsSessionRepository) *CreateAndSendMsgToEveryoneExcept {
	return &CreateAndSendMsgToEveryoneExcept{sessionsRepo: sessionsRepo}
}

func (r *CreateAndSendMsgToEveryoneExcept) CreateAndSendMsgToEveryoneExcept(msg model.CreateMessage, ids []model.ID) error {
	newMsg := model.Message{
		Id:      model.NewID(uuid.NewString()),
		ReplyId: msg.ReplyId,
		Date:    msg.Date,
		Author:  msg.Author,
		From:    msg.From,
		Silent:  msg.Silent,
		Body:    msg.Body,
	}
	connected, err := r.sessionsRepo.Connected()
	if err != nil {
		return fmt.Errorf("connected clients: %s", err)
	}

	for _, client := range connected {
		var isExcepted = slices.ContainsFunc(ids, func(exceptedId model.ID) bool {
			return client.Id == exceptedId
		})
		if !isExcepted {
			r.sessionsRepo.Send(newMsg, client.Id)
		}
	}

	return nil
}

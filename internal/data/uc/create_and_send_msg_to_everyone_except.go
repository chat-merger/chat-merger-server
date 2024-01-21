package uc

import (
	"chatmerger/internal/component/eventbus"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/domain/repository"
	"chatmerger/internal/usecase"
	"fmt"
	"github.com/google/uuid"
	"slices"
)

var _ usecase.CreateAndSendMsgToEveryoneExceptUc = (*CreateAndSendMsgToEveryoneExcept)(nil)

type CreateAndSendMsgToEveryoneExcept struct {
	cRepo repository.ClientsRepository
	bus   *eventbus.EventBus
}

func NewCreateAndSendMsgToEveryoneExcept(
	cRepo repository.ClientsRepository,
	bus *eventbus.EventBus,
) *CreateAndSendMsgToEveryoneExcept {
	return &CreateAndSendMsgToEveryoneExcept{bus: bus, cRepo: cRepo}
}

func (r *CreateAndSendMsgToEveryoneExcept) CreateAndSendMsgToEveryoneExcept(msg model.CreateMessage, ids ...model.ID) (*model.Message, error) {
	newMsg := &model.Message{
		Id:       model.NewID(uuid.NewString()),
		ReplyId:  msg.ReplyId,
		Date:     msg.Date,
		Username: msg.Username,
		From:     msg.From,
		Silent:   msg.Silent,
		Body:     msg.Body,
	}

	expStatus := model.ConnStatusActive
	connected, err := r.cRepo.GetClients(model.ClientsFilter{Status: &expStatus})
	if err != nil {
		return nil, fmt.Errorf("get clients: %s", err)
	}

	// definition client who received msg
	recipients := make([]model.ID, 0, len(connected))
	for _, client := range connected {
		var isExcepted = slices.ContainsFunc(ids, func(exceptedId model.ID) bool {
			return client.Id == exceptedId
		})
		if !isExcepted {
			recipients = append(recipients, client.Id)
		}
	}

	err = r.bus.Publish(eventbus.Event{Message: newMsg}, recipients...)
	if err != nil {
		return nil, fmt.Errorf("publish new message: %s", err)
	}

	return newMsg, nil
}

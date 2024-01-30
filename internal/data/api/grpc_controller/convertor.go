package grpc_controller

import (
	"chatmerger/internal/data/api/pb"
	"chatmerger/internal/domain/model"
	"time"
)

func newMsgToDomain(request *pb.NewMessageBody, client string) (*model.CreateMessage, error) {

	msg := &model.CreateMessage{
		ReplyId:   (*model.ID)(request.ReplyMsgId),
		Date:      time.Unix(request.CreatedAt, 0),
		Username:  request.Username,
		From:      client,
		Silent:    request.Silent,
		Text:      request.Text,
		Media:     make([]model.Media, 0, len(request.Media)),
		Forwarded: make([]model.Forward, 0, len(request.Forwarded)),
	}
	// add media
	for _, it := range request.Media {
		msg.Media = append(msg.Media, mediaToDomain(it))
	}
	for _, it := range request.Forwarded {
		msg.Forwarded = append(msg.Forwarded, forwardToDomain(it))
	}

	return msg, nil
}

func newMsgToPb(msg model.Message) (*pb.Message, error) {
	var replyMsgId *string
	if msg.ReplyId != nil {
		replyMsgId = (*string)(msg.ReplyId)
	}
	// response
	response := &pb.Message{
		Id:         string(msg.Id),
		Client:     msg.From,
		CreatedAt:  msg.Date.Unix(),
		Silent:     msg.Silent,
		ReplyMsgId: replyMsgId,
		Username:   msg.Username,
		Text:       msg.Text,
		Media:      make([]*pb.Media, 0, len(msg.Media)),
		Forwarded:  make([]*pb.Forwarded, 0, len(msg.Forwarded)),
	}
	// add media
	for _, it := range msg.Media {
		response.Media = append(response.Media, mediaToPb(it))
	}
	// add forwarded
	for _, it := range msg.Forwarded {
		response.Forwarded = append(response.Forwarded, forwardToPb(it))
	}

	return response, nil
}

func forwardToPb(bm model.Forward) *pb.Forwarded {
	media := make([]*pb.Media, 0, len(bm.Media))
	for _, it := range bm.Media {
		media = append(media, mediaToPb(it))
	}
	return &pb.Forwarded{
		Id:        (*string)(bm.Id),
		CreatedAt: bm.Date.Unix(),
		Username:  bm.Username,
		Text:      bm.Text,
		Media:     media,
	}
}

func forwardToDomain(bm *pb.Forwarded) model.Forward {
	media := make([]model.Media, 0, len(bm.Media))
	for _, it := range bm.Media {
		media = append(media, mediaToDomain(it))
	}
	return model.Forward{
		Id:       (*model.ID)(bm.Id),
		Date:     time.Unix(bm.CreatedAt, 0),
		Username: bm.Username,
		Text:     bm.Text,
		Media:    media,
	}
}

func mediaToDomain(bm *pb.Media) model.Media {
	return model.Media{
		Kind:    mediaTypeToDomain(bm.Type),
		Spoiler: bm.Spoiler,
		Url:     bm.Url,
	}
}

func mediaToPb(bm model.Media) *pb.Media {
	return &pb.Media{
		Type:    mediaTypeToPb(bm.Kind),
		Spoiler: bm.Spoiler,
		Url:     bm.Url,
	}
}

func mediaTypeToDomain(kind pb.Media_Type) model.MediaType {
	var tf model.MediaType
	switch kind {
	case pb.Media_AUDIO:
		tf = model.Audio
	case pb.Media_VIDEO:
		tf = model.Video
	case pb.Media_FILE:
		tf = model.File
	case pb.Media_PHOTO:
		tf = model.Photo
	case pb.Media_STICKER:
		tf = model.Sticker
	}
	return tf
}

func mediaTypeToPb(kind model.MediaType) pb.Media_Type {
	var tf pb.Media_Type
	switch kind {
	case model.Audio:
		tf = pb.Media_AUDIO
	case model.Video:
		tf = pb.Media_VIDEO
	case model.File:
		tf = pb.Media_FILE
	case model.Photo:
		tf = pb.Media_PHOTO
	case model.Sticker:
		tf = pb.Media_STICKER
	}
	return tf
}

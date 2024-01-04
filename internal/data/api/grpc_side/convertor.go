package grpc_side

import (
	"chatmerger/internal/data/api/pb"
	"chatmerger/internal/domain/model"
	"errors"
	"log"
	"time"
)

func requestToCreateMessage(request *pb.Request, client string) (*model.CreateMessage, error) {
	var body model.Body
	switch request.Body.(type) {
	case *pb.Request_Text:
		var rt = request.Body.(*pb.Request_Text).Text
		body = &model.BodyText{
			Format: pbTextFormatToModel(rt.Format),
			Value:  rt.Value,
		}
	case *pb.Request_Media:
		var rm = request.Body.(*pb.Request_Media).Media
		body = &model.BodyMedia{
			Kind:    pbMediaTypeToModel(rm.Type),
			Caption: rm.Caption,
			Spoiler: rm.Spoiler,
			Url:     rm.Url,
		}
	default:
		return nil, errors.New("request body not match with RequestBody interface")
	}

	return &model.CreateMessage{
		ReplyId: stringToIdIfExists(request.ReplyMsgId),
		Date:    time.Unix(request.CreatedAt, 0),
		Author:  request.Author,
		From:    client,
		Silent:  request.IsSilent,
		Body:    body,
	}, nil
}

func messageToResponse(msg model.Message) (*pb.Response, error) {
	var replyMsgId *string
	if msg.ReplyId != nil {
		id := msg.ReplyId.Value()
		replyMsgId = &id
	}
	// response
	response := &pb.Response{
		Id:         msg.Id.Value(),
		ReplyMsgId: replyMsgId,
		CreatedAt:  msg.Date.Unix(),
		Author:     msg.Author,
		Client:     msg.From,
		IsSilent:   msg.Silent,
		Body:       nil, // WithoutBody!!!!!
	}
	// add body
	switch msg.Body.(type) {
	case *model.BodyText:
		text := msg.Body.(*model.BodyText)
		response.Body = modelBodyTextToPb(*text)
	case *model.BodyMedia:
		media := msg.Body.(*model.BodyMedia)
		response.Body = modelBodyMediaToPb(*media)
	default:
		log.Fatalf("unknown msg.Body:  %#v", msg.Body)
	}
	return response, nil
}

func modelBodyTextToPb(bt model.BodyText) *pb.Response_Text {
	return &pb.Response_Text{
		Text: &pb.Text{
			Format: modelTextFormatToPbTextFormat(bt.Format),
			Value:  bt.Value,
		},
	}
}

func modelBodyMediaToPb(bm model.BodyMedia) *pb.Response_Media {
	return &pb.Response_Media{
		Media: &pb.Media{
			Type:    modelMediaTypeToPbMediaType(bm.Kind),
			Caption: bm.Caption,
			Spoiler: bm.Spoiler,
			Url:     bm.Url,
		},
	}
}

func stringToIdIfExists(str *string) *model.ID {
	if str == nil {
		return nil
	}
	var id = model.NewID(*str)
	return &id
}

func pbTextFormatToModel(format pb.Text_Format) model.TextFormat {
	var tf model.TextFormat
	switch format {
	case pb.Text_MARKDOWN:
		tf = model.Markdown
	case pb.Text_PLAIN:
		tf = model.Plain
	}
	return tf
}

func pbMediaTypeToModel(kind pb.Media_Type) model.MediaType {
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

func modelTextFormatToPbTextFormat(format model.TextFormat) pb.Text_Format {
	var tf pb.Text_Format
	switch format {
	case model.Markdown:
		tf = pb.Text_MARKDOWN
	case model.Plain:
		tf = pb.Text_PLAIN
	}
	return tf
}

func modelMediaTypeToPbMediaType(kind model.MediaType) pb.Media_Type {
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

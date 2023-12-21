package pb

import (
	"chatmerger/internal/domain/model"
	"errors"
	"log"
	"time"
)

func RequestToCreateMessage(request *Request, client string) (*model.CreateMessage, error) {
	var body model.Body
	switch request.Body.(type) {
	case *Request_Text:
		var rt = request.Body.(*Request_Text).Text
		body = &model.BodyText{
			Format: PbTextFormatToModel(rt.Format),
			Value:  rt.Value,
		}
	case *Request_Media:
		var rm = request.Body.(*Request_Media).Media
		body = &model.BodyMedia{
			Kind:    PbMediaTypeToModel(rm.Type),
			Caption: rm.Caption,
			Spoiler: rm.Spoiler,
			Url:     rm.Url,
		}
	default:
		return nil, errors.New("request body not match with RequestBody interface")
	}

	return &model.CreateMessage{
		ReplyId: StringToIdIfExists(request.ReplyMsgId),
		Date:    time.Unix(request.CreatedAt, 0),
		Author:  request.Author,
		From:    client,
		Silent:  request.IsSilent,
		Body:    body,
	}, nil
}

func MessageToResponse(msg model.Message) (*Response, error) {
	var replyMsgId *string
	if msg.ReplyId != nil {
		id := msg.ReplyId.Value()
		replyMsgId = &id
	}
	// response
	response := &Response{
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
		response.Body = ModelBodyTextToPb(*text)
	case *model.BodyMedia:
		media := msg.Body.(*model.BodyMedia)
		response.Body = ModelBodyMediaToPb(*media)
	default:
		log.Fatalf("unknown msg.Body:  %#v", msg.Body)
	}
	return response, nil
}

func ModelBodyTextToPb(bt model.BodyText) *Response_Text {
	return &Response_Text{
		Text: &Text{
			Format: ModelTextFormatToPbTextFormat(bt.Format),
			Value:  bt.Value,
		},
	}
}

func ModelBodyMediaToPb(bm model.BodyMedia) *Response_Media {
	return &Response_Media{
		Media: &Media{
			Type:    ModelMediaTypeToPbMediaType(bm.Kind),
			Caption: bm.Caption,
			Spoiler: bm.Spoiler,
			Url:     bm.Url,
		},
	}
}

func StringToIdIfExists(str *string) *model.ID {
	if str == nil {
		return nil
	}
	var id = model.NewID(*str)
	return &id
}

func PbTextFormatToModel(format Text_Format) model.TextFormat {
	var tf model.TextFormat
	switch format {
	case Text_MARKDOWN:
		tf = model.Markdown
	case Text_PLAIN:
		tf = model.Plain
	}
	return tf
}

func PbMediaTypeToModel(kind Media_Type) model.MediaType {
	var tf model.MediaType
	switch kind {
	case Media_AUDIO:
		tf = model.Audio
	case Media_VIDEO:
		tf = model.Video
	case Media_FILE:
		tf = model.File
	case Media_PHOTO:
		tf = model.Photo
	case Media_STICKER:
		tf = model.Sticker
	}
	return tf
}

func ModelTextFormatToPbTextFormat(format model.TextFormat) Text_Format {
	var tf Text_Format
	switch format {
	case model.Markdown:
		tf = Text_MARKDOWN
	case model.Plain:
		tf = Text_PLAIN
	}
	return tf
}

func ModelMediaTypeToPbMediaType(kind model.MediaType) Media_Type {
	var tf Media_Type
	switch kind {
	case model.Audio:
		tf = Media_AUDIO
	case model.Video:
		tf = Media_VIDEO
	case model.File:
		tf = Media_FILE
	case model.Photo:
		tf = Media_PHOTO
	case model.Sticker:
		tf = Media_STICKER
	}
	return tf
}

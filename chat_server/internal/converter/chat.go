package converter

import (
	"chat/chat_server/internal/model"
	desc "chat/chat_server/pkg/chat_v1"
)

func ToChatCreateFromDesc(req *desc.CreateRequest) *model.ChatCreate {
	return &model.ChatCreate{
		Usernames: req.GetUsernames(),
	}
}

func ToMessageFromDesc(req *desc.SendMessageRequest) *model.Message {
	return &model.Message{
		From:      req.GetFrom(),
		Text:      req.GetText(),
		Timestamp: req.GetTimestamp().AsTime(),
	}
}

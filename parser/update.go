package parser

import (
	"encoding/json"
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
	"telegram-parser/db"
)

// GetUpdates starts reading messages from the telegram raw updates channel.
func (a *App) GetUpdates() {
	var COUNT_REVIEW_GOROUTINES = 100
	var UPDATES_CHANNEL_CAPACITY = 100000

	// rawUpdates gets all updates comming from tdlib
	rawUpdates := a.Telegram.Client.GetRawUpdatesChannel(UPDATES_CHANNEL_CAPACITY)

	for i := 0; i < COUNT_REVIEW_GOROUTINES; i++ {
		go a.messageReview(rawUpdates)
	}

	logrus.Info("The system for processing Telegram updates was launched.")
	logrus.Infof("%v goroutine messageReview launched.", COUNT_REVIEW_GOROUTINES)
}

// messageReview checks the message for compliance with system requirements.
func (a *App) messageReview(rawUpdates <-chan tdlib.UpdateMsg) {
	var (
		tgCli  = a.Telegram.Client
		dbConn = a.DbCli
		rabbit = a.MqCli
	)

	for update := range rawUpdates {

		logrus.Info("GETTED UPDATE")

		if update.Data["@type"] != "updateChatLastMessage" {
			logrus.Warnf("Received update is of type `%v`, type `updateChatLastMessage` is required.", update.Data["@type"])
			continue
		}

		var lastMsgUpdate *tdlib.UpdateChatLastMessage
		err := json.Unmarshal(update.Raw, &lastMsgUpdate)
		if err != nil {
			// TODO убрать fatal
			logrus.Fatal(err)
		}

		if lastMsgUpdate.LastMessage == nil {
			logrus.Error("lastMsgUpdate have not LastMessage")
			continue
		}

		if lastMsgUpdate.LastMessage.Content.GetMessageContentEnum() != "messageText" {
			logrus.Info("Content is not text")
			continue
		}

		chat, err := tgCli.GetChat(lastMsgUpdate.ChatID)
		if err != nil {
			logrus.Panicf("%v with lastMsgUpdate= %v", err, lastMsgUpdate)
		}

		// Убрано на время тестирования
		if chat.Type.GetChatTypeEnum() != "chatTypeSupergroup" {
			logrus.Warnf("Chat %v is not super group. It is %v\n", chat.Title, chat.Type.GetChatTypeEnum())
			continue
		}

		link, err := tgCli.GetMessageLink(chat.ID, lastMsgUpdate.LastMessage.ID, true, true)
		if err != nil {
			logrus.Panic("Link: ", err)
		}

		m := db.NewMessage(lastMsgUpdate.LastMessage, chat.Title, link)

		err = dbConn.InsertMessage(m)
		if err != nil {
			logrus.Error(err)
		}

		err = rabbit.Publish(m)
		if err != nil {
			logrus.Panic(err)
		}

		logrus.Infof("Сообщение с id %v передано в rabbit\n", m.MessageID)
	}
}

// extra

///*GetChatList Returns an ordered list of chats in a chat list. Chats are sorted by the pair (chat.position.order, chat.id) in descending order. @param limit The maximum number of chats to be returned. It is possible that fewer chats than the limit are returned even if the end of the list is not reached
// */
//func (t *Telegram) GetChatList(limit int) ([]*tdlib.Chat, error) {
//	var allChats []*tdlib.Chat
//	var offsetOrder = int64(math.MaxInt64)
//	var offsetChatID = int64(0)
//	var chatList = tdlib.NewChatListMain()
//	var lastChat *tdlib.Chat
//
//	for len(allChats) < limit {
//		if len(allChats) > 0 {
//			lastChat = allChats[len(allChats)-1]
//			for i := 0; i < len(lastChat.Positions); i++ {
//				//Find the main chat list
//				if lastChat.Positions[i].List.GetChatListEnum() == tdlib.ChatListMainType {
//					offsetOrder = int64(lastChat.Positions[i].Order)
//				}
//			}
//			offsetChatID = lastChat.ID
//		}
//
//		var chats, err = t.Client.GetChats(chatList, tdlib.JSONInt64(offsetOrder),
//			offsetChatID, int32(limit-len(allChats)))
//		if err != nil {
//			return nil, err
//		}
//		if len(chats.ChatIDs) == 0 {
//			return allChats, nil
//		}
//
//		for _, chatID := range chats.ChatIDs {
//			chat, err := t.Client.GetChat(chatID)
//			if err != nil {
//				return nil, err
//			}
//			allChats = append(allChats, chat)
//		}
//	}
//	return allChats, nil
//}

package app

import (
	"encoding/json"
	"fmt"
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
	"os"
	"telegram-parser/db"
)

//    --------------------------------------------------------------------------------
//                                    STRUCTS
//    --------------------------------------------------------------------------------

// TODO change chan to rabbitmq
type Telegram struct {
	Client          *tdlib.Client    // Telegram.org client
	ReceivedUpdates chan *db.Message // Contains messages for further processing (distribution by nodes to track statistics)
}

//    --------------------------------------------------------------------------------
//                                     METHODS
//    --------------------------------------------------------------------------------

// NewTgClient Create new instance of client
func NewTgClient() *Telegram {
	client := tdlib.NewClient(tdlib.Config{
		APIID:              os.Getenv("TGAPIID"),
		APIHash:            os.Getenv("TGAPIHASH"),
		SystemLanguageCode: "en",
		DeviceModel:        "Server",
		SystemVersion:      "1.0.0",
		ApplicationVersion: "1.0.0",
		// Optional fields
		UseTestDataCenter:      false,
		DatabaseDirectory:      "./tdlib-db",
		FileDirectory:          "./tdlib-files",
		UseFileDatabase:        false,
		UseChatInfoDatabase:    false,
		UseMessageDatabase:     false,
		UseSecretChats:         false,
		EnableStorageOptimizer: false,
		IgnoreFileNames:        false,
	})
	updates := make(chan *db.Message, 1000)

	return &Telegram{Client: client, ReceivedUpdates: updates}
}

// Authorization is used to authorize the user
func (t *Telegram) Authorization() {
	for {
		currentState, _ := t.Client.Authorize()
		switch currentState.GetAuthorizationStateEnum() {
		case tdlib.AuthorizationStateWaitPhoneNumberType:
			_, err := t.Client.SendPhoneNumber(os.Getenv("TGTELEPHONENUMBER"))
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		case tdlib.AuthorizationStateWaitCodeType:
			var code string
			fmt.Print("Enter code: ")
			fmt.Scanln(&code)
			_, err := t.Client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		case tdlib.AuthorizationStateWaitPasswordType:
			var password string
			fmt.Print("Enter Password: ")
			fmt.Scanln(&password)
			_, err := t.Client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		case tdlib.AuthorizationStateReadyType:
			fmt.Println("Authorization Ready.\n")
			return
		}
	}
}

func (t *Telegram) RunHandlingUpdates() {
	postgresClient, err := db.ConnectToPostgres()
	if err != nil {
		//	TODO что делать с ошибкой (без базы невозможно парсить сообщения)
		panic(err)
	}
	go t.MessagesHandling(postgresClient)
	go t.GetUpdates()
}

// GetUpdates catches records only about new unread messages in channels
func (t *Telegram) GetUpdates() {
	// rawUpdates gets all updates comming from tdlib
	rawUpdates := t.Client.GetRawUpdatesChannel(100)
	for update := range rawUpdates {

		var updateLastMessage tdlib.UpdateChatLastMessage
		err := json.Unmarshal(update.Raw, &updateLastMessage)
		if err != nil {
			//	TODO придумать что делать с этой ошибкой (можем потерять сообщения)
			logrus.Panic(err)
		}

		if updateLastMessage.Type != "updateChatLastMessage" {
			continue
		}

		chat, err := t.Client.GetChat(updateLastMessage.ChatID)
		if err != nil {
			//	TODO придумать что делать с этой ошибкой (можем потерять сообщения)
			logrus.Panic(err)
		}

		if chat.Type.GetChatTypeEnum() != "chatTypeSupergroup" {
			continue
		}

		if updateLastMessage.LastMessage.Content.GetMessageContentEnum() != "messageText" {
			continue
		}

		m := db.NewMessage(updateLastMessage.LastMessage, chat)

		t.ReceivedUpdates <- m
	}
}

func (t *Telegram) MessagesHandling(dbClient db.DB) {
	for update := range t.ReceivedUpdates {
		// Add a new message to the database
		err := dbClient.Insert(update)
		if err != nil {
			//	TODO придумать что делать с этой ошибкой (можем потерять сообщения)
			logrus.Panic(err)
		}

		//	TODO сообщение отправляется в сервис (с помощью консистентного хеширования) для дальнейшего наблюдения

	}
}

// TODO updateLastMessage.LastMessage - положить в базу и направить на дальнейшую обработку

//    --------------------------------------------------------------------------------
//                                        EXTRA
//    --------------------------------------------------------------------------------

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

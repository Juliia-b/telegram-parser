package parser

import (
	"encoding/json"
	"fmt"
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"telegram-parser/db"
	"telegram-parser/mq"
)

/*---------------------------------STRUCTURES----------------------------------------*/

type App struct {
	Telegram *Telegram
	DbCli    db.DB
	MqCli    *mq.Rabbit
}

type Telegram struct {
	Client *tdlib.Client // Telegram.org client
}

/*-----------------------------------METHODS-----------------------------------------*/

// AppInstance returns a structure with connections to all services
func AppInstance(dbClient db.DB, mqClient *mq.Rabbit) *App {
	tg := newTgClient()

	return &App{
		Telegram: tg,
		DbCli:    dbClient,
		MqCli:    mqClient,
	}
}

// newTgClient Create new instance of client
func newTgClient() *Telegram {
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

	return &Telegram{Client: client}
}

// TelegramAuthorization is used to authorize the user
func (a *App) TelegramAuthorization(wg *sync.WaitGroup) {
	tgCli := a.Telegram.Client

	for {
		currentState, _ := tgCli.Authorize()
		switch currentState.GetAuthorizationStateEnum() {
		case tdlib.AuthorizationStateWaitPhoneNumberType:
			_, err := tgCli.SendPhoneNumber(os.Getenv("TGTELNUMBER"))
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		case tdlib.AuthorizationStateWaitCodeType:
			var code string
			fmt.Print("Enter code: ")
			fmt.Scanln(&code)
			_, err := tgCli.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		case tdlib.AuthorizationStateWaitPasswordType:
			var password string
			fmt.Print("Enter Password: ")
			fmt.Scanln(&password)
			_, err := tgCli.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		case tdlib.AuthorizationStateReadyType:
			logrus.Info("TelegramAuthorization Ready.\n")
			wg.Done()
			return
		}
	}
}

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

		logrus.Infof("GET NEW messageReview")

		var updateLastMessage tdlib.UpdateChatLastMessage
		err := json.Unmarshal(update.Raw, &updateLastMessage)
		if err != nil {
			logrus.Panic(err)
		}

		if updateLastMessage.Type != "updateChatLastMessage" {
			//logrus.Info("Not type updateChatLastMessage")
			continue
		}

		if updateLastMessage.LastMessage == nil {
			logrus.Error("updateLastMessage have not LastMessage")
			continue
		}

		if updateLastMessage.LastMessage.Content.GetMessageContentEnum() != "messageText" {
			logrus.Info("Content is not text")
			continue
		}

		chat, err := tgCli.GetChat(updateLastMessage.ChatID)
		if err != nil {
			logrus.Panicf("%v with updateLastMessage= %v", err, updateLastMessage)
		}

		// Убрано на время тестирования
		if chat.Type.GetChatTypeEnum() != "chatTypeSupergroup" {
			logrus.Warnf("Chat %v is not super group. It is %v\n", chat.Title, chat.Type.GetChatTypeEnum())
			continue
		}

		m := db.NewMessage(updateLastMessage.LastMessage, chat.Title)

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

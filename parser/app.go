package parser

import (
	"encoding/json"
	"fmt"
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
	"os"
	"telegram-parser/db"
	"telegram-parser/mq"
)

//    --------------------------------------------------------------------------------
//                                  STRUCTURES
//    --------------------------------------------------------------------------------

type App struct {
	Telegram *Telegram
	DbCli    db.DB
	MqCli    *mq.Rabbit
}

type Telegram struct {
	Client *tdlib.Client // Telegram.org client
}

//    --------------------------------------------------------------------------------
//                                     METHODS
//    --------------------------------------------------------------------------------

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
func (a *App) TelegramAuthorization() {
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

		var updateLastMessage tdlib.UpdateChatLastMessage
		err := json.Unmarshal(update.Raw, &updateLastMessage)
		if err != nil {
			logrus.Panic(err)
		}

		if updateLastMessage.Type != "updateChatLastMessage" {
			continue
		}

		if !updateLastMessage.LastMessage.CanGetStatistics {
			continue
		}

		chat, err := tgCli.GetChat(updateLastMessage.ChatID)
		if err != nil {
			logrus.Panicf("%v with updateLastMessage= %v", err, updateLastMessage)
		}

		// Убрано на время тестирования
		//if chat.Type.GetChatTypeEnum() != "chatTypeSupergroup" {
		//	continue
		//}

		if updateLastMessage.LastMessage.Content.GetMessageContentEnum() != "messageText" {
			continue
		}

		m := db.NewMessage(updateLastMessage.LastMessage, chat.Title)

		err = dbConn.Insert(m)
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

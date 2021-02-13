package parser

import (
	"encoding/json"
	"fmt"
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
	"telegram-parser/db"
	"telegram-parser/flags"
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
func AppInstance(config *flags.Config) *App {
	tg := newTgClient(config)

	dbClient, err := db.ConnectToPostgres(config)
	if err != nil {
		logrus.Panic(err)
	}

	mq, err := mq.RabbitInit()
	if err != nil {
		logrus.Panic(err)
	}

	return &App{
		Telegram: tg,
		DbCli:    dbClient,
		MqCli:    mq,
	}
}

// newTgClient Create new instance of client
func newTgClient(conf *flags.Config) *Telegram {
	client := tdlib.NewClient(tdlib.Config{
		APIID:              conf.Telegram.APIID,
		APIHash:            conf.Telegram.APIHash,
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
func (a *App) TelegramAuthorization(conf *flags.Config) {
	tgCli := a.Telegram.Client

	for {
		currentState, _ := tgCli.Authorize()
		switch currentState.GetAuthorizationStateEnum() {
		case tdlib.AuthorizationStateWaitPhoneNumberType:
			_, err := tgCli.SendPhoneNumber(conf.Telegram.TelephoneNumber)
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

// GetUpdates catches records only about new unread messages in channels
func (a *App) GetUpdates() {
	logrus.Info("RUN GETTING UPDATES")

	var (
		tgCli  = a.Telegram.Client
		dbconn = a.DbCli
		rabbit = a.MqCli
	)

	// rawUpdates gets all updates comming from tdlib
	rawUpdates := tgCli.GetRawUpdatesChannel(100000)
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

		err = dbconn.Insert(m)
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

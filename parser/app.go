package parser

import (
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

package main

import (
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
	"sync"
	"telegram-parser/db"
	"telegram-parser/helpers"
	"telegram-parser/mq"
	"telegram-parser/parser"
	"telegram-parser/server"
)

func main() {
	//System Setup
	helpers.ConfigureLogrus()
	helpers.CheckEnv()

	dbClient := db.ConnectToPostgres()
	mqClient := mq.RabbitInit()

	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")

	app := parser.AppInstance(dbClient, mqClient)

	var wg sync.WaitGroup
	wg.Add(1)
	app.TelegramAuthorization(&wg)
	wg.Wait()

	// Необходимо, чтобы tdlib знал о чатах, поэтому сначала нужно пройтись по всем чатам
	_, err := app.Telegram.GetChatList(5000)
	if err != nil {
		logrus.Panicf("Fail to get chat list with error = '%v'.", err.Error())
	}

	// Run handling updates from Telegram
	go app.GetUpdates()

	// Run tracking statistics
	app.StartTrackingStatistics(10) // TODO в продакшене заменить на большее число

	s := server.Init(dbClient)

	logrus.Info("Server is running on ", s.Server.Addr)
	logrus.Panic(s.Server.ListenAndServe())
}

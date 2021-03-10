package main

import (
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
	"telegram-parser/db"
	"telegram-parser/helpers"
	"telegram-parser/mq"
	"telegram-parser/parser"
	"telegram-parser/server"
	"time"

	"net/http"
	_ "net/http/pprof"
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
	app.TelegramAuthorization()

	// It is necessary that tdlib knows about chats, so first you need to go through all chats
	if _, err := app.Telegram.GetChatList(1000000); err != nil {
		logrus.Panicf("Fail to get chat list with error = '%v'.", err.Error())
	}

	time.Sleep(10 * time.Second)

	// Run handling updates from Telegram
	go app.GetUpdates()

	// Run tracking statistics
	app.StartTrackingStatistics(10)

	// Run server for pprof
	go func() {
		logrus.Fatal(http.ListenAndServe(":3000", nil))
	}()

	s := server.Init(dbClient)

	logrus.Info("Server is running on ", s.Server.Addr)
	logrus.Panic(s.Server.ListenAndServe())
}

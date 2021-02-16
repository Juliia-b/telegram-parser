package main

import (
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
	"telegram-parser/db"
	"telegram-parser/handler"
	"telegram-parser/helpers"
	"telegram-parser/mq"
	"telegram-parser/parser"
)

func main() {
	//System Setup
	helpers.ConfigureLogrus()
	helpers.CheckEnv()

	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")
	//
	dbClient, err := db.ConnectToPostgres()
	if err != nil {
		logrus.Panic(err)
	}

	mqClient, err := mq.RabbitInit()
	if err != nil {
		logrus.Panic(err)
	}

	app := parser.AppInstance(dbClient, mqClient)
	app.TelegramAuthorization()

	//chats, err := app.Telegram.GetChatList(5000)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, chat := range chats {
	//	logrus.Infof("TITLE: %v\nCHAT TYPE: %v\nCAN GET STATISTICS: %v\n----------------------------\n", chat.Title, chat.Type.GetChatTypeEnum(), chat.LastMessage.CanGetStatistics)
	//}

	// TODO убрать из списка выдачи лучших записей все у которых значения стоят на 0 (просмотры) или 1

	// Run handling updates from Telegram
	go app.GetUpdates()

	// Run
	app.StartTrackingStatistics(50)

	// -------------

	r := handler.RouterInit(dbClient)

	logrus.Info("Server is running on ", r.Server.Addr)
	logrus.Panic(r.Server.ListenAndServe())

	forever := make(chan bool)
	<-forever
}

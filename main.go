package main

import (
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
	"sync"
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

	dbClient, err := db.ConnectToPostgres()
	if err != nil {
		logrus.Panic(err)
	}

	mqClient, err := mq.RabbitInit()
	if err != nil {
		logrus.Panic(err)
	}

	app := parser.AppInstance(dbClient, mqClient)

	var wg sync.WaitGroup
	wg.Add(1)
	app.TelegramAuthorization(&wg)
	wg.Wait()

	// Run handling updates from Telegram
	go app.GetUpdates()

	// Run tracking statistics
	app.StartTrackingStatistics(50)

	r := handler.ServerInit(dbClient)

	logrus.Info("Server is running on ", r.Server.Addr)
	logrus.Panic(r.Server.ListenAndServe())
}

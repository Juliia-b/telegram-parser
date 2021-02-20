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

	// Run handling updates from Telegram
	go app.GetUpdates()

	// Run tracking statistics
	app.StartTrackingStatistics(50)

	r := server.Init(dbClient)

	logrus.Info("Server is running on ", r.Server.Addr)
	logrus.Panic(r.Server.ListenAndServe())
}

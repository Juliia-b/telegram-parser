package main

import (
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
	"telegram-parser/app"
	"telegram-parser/db"
	"telegram-parser/flags"
	"telegram-parser/grpc_client"
	"telegram-parser/handling"
	"telegram-parser/helpers"
	"telegram-parser/mq"
)

func main() {
	// System Setup
	helpers.CheckENV()
	helpers.ConfigureLogrus()

	// инициализация канала для новых адресов нод
	serviceAddresses := make(chan string)
	// инициализация структуры подключения к клиентам сервиса парсинга
	serviceConnections := grpc_client.NewParserConnections()

	// база
	postgresClient, err := db.ConnectToPostgres()
	if err != nil {
		logrus.Panic(err)
	}

	// очередь сообщений
	mq, err := mq.RabbitInit()
	if err != nil {
		logrus.Panic(err)
	}

	// структура общего назначения с доступом к базе, брокеру сообщений, подключения к клиентам сервиса и информации о состоянии консистентного хеша
	handle := handling.NewHandleStruct(serviceConnections, postgresClient, mq)

	// пытается подключиться ко всем адресам в списке
	go handle.HandleNewServiceAddresses(serviceAddresses)
	// парсит флаги. в случае если аргумент является ipv4 адресом он отправляется в канал serviceAddresses
	flags.ParseFlags(serviceAddresses)

	// запуск консьюмера очереди для получения сообщений и их обработки ( нахождение адреса ноды для отправки, отправка)
	handle.RunHandlingMsgsFromMQ(20)

	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")

	telegramCli := app.NewTgClient()
	telegramCli.Authorization()

	go telegramCli.GetUpdates(handle.Rabbit)

	forever := make(chan bool)
	<-forever
}

package main

import (
	"github.com/Arman92/go-tdlib"
	"github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
	"os"
	"telegram-parser/app"
	"telegram-parser/db"
	"telegram-parser/flags"
	"telegram-parser/grpc_client"
	"telegram-parser/handling"
	"telegram-parser/mq"
)

func main() {
	// System Setup
	checkENV()
	configureLogrus()

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
	handle := handling.NewHandleStruct(postgresClient, mq, serviceConnections)

	// пытается подключиться ко всем адресам в списке
	go handle.HandleNewServiceAddresses(serviceAddresses)
	// парсит флаги. в случае если аргумент является ipv4 адресом он отправляется в канал serviceAddresses
	flags.ParseFlags(serviceAddresses)

	// запуск консьюмера очереди для получения сообщений и их обработки ( нахождение адреса ноды для отправки, отправка)
	go handle.HandleMsgsFromMQ()

	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")

	telegramCli := app.NewTgClient()
	telegramCli.Authorization()

	go telegramCli.GetUpdates(handle.Rabbit)

	// ----------------------------------------
	//
	//
	//
	//
	//
	//
	//conf := flags.ParseFlags()
	//
	//parserConn := grpc_client.NewParserConnections()
	//parserConn.ParserClientsInit(conf.ParserAddrs)

	//tdlib.SetLogVerbosityLevel(1)
	//tdlib.SetFilePath("./errors.txt")
	//
	//telegram := app.NewTgClient()
	//telegram.Authorization()
	//
	//telegram.RunHandlingUpdates()

	//app.RedBlackTree()

	forever := make(chan bool)
	<-forever
}

func checkENV() {
	if os.Getenv("POSTGRESPASSWORD") == "" {
		logrus.Panic(`Environment variable "POSTGRESPASSWORD" is blank`)
	}

	if os.Getenv("TGTELEPHONENUMBER") == "" {
		logrus.Panic(`Environment variable "TGTELEPHONENUMBER" is blank`)
	}

	if os.Getenv("TGAPIID") == "" {
		logrus.Panic(`Environment variable "TGAPIID" is blank`)
	}

	if os.Getenv("TGAPIHASH") == "" {
		logrus.Panic(`Environment variable "TGAPIHASH" is blank`)
	}
}

func configureLogrus() {
	formatter := runtime.Formatter{
		ChildFormatter: &logrus.TextFormatter{
			TimestampFormat: "02-01-2006 15:04:05", // "Mon Jan 2 15:04:05 MST 2006"
			FullTimestamp:   true,
		},
		File:         true,
		BaseNameOnly: true}

	formatter.Line = true

	logrus.SetFormatter(&formatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

package main

import (
	"github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {

	configureLogrus()

	// ----------------------------
	//checkENV()
	//configureLogrus()
	//
	//flags.ParseFlags()

	//tdlib.SetLogVerbosityLevel(1)
	//tdlib.SetFilePath("./errors.txt")
	//
	//telegram := app.NewTgClient()
	//telegram.Authorization()
	//
	//telegram.RunHandlingUpdates()
	//
	//cli, _ := db.ConnectToPostgres()
	////cli.GetAllData()
	//cli.GetMessageById(-1001129804770, 408447614976)

	for {
	}
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

package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"telegram-parser/flags"
)

func main() {
	checkENV()
	configureLogrus()

	//flags.ParseFlags()
	flags.ParseFlags()

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
		panic(`Environment variable "POSTGRESPASSWORD" is blank`)
	}

	if os.Getenv("TGTELEPHONENUMBER") == "" {
		panic(`Environment variable "TGTELEPHONENUMBER" is blank`)
	}

	if os.Getenv("TGAPIID") == "" {
		panic(`Environment variable "TGAPIID" is blank`)
	}

	if os.Getenv("TGAPIHASH") == "" {
		panic(`Environment variable "TGAPIHASH" is blank`)
	}
}

func configureLogrus() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05" // "Mon Jan 2 15:04:05 MST 2006"
	Formatter.FullTimestamp = true

	log.SetFormatter(Formatter)
}

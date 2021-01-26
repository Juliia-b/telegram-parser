package main

import (
	"fmt"
	"github.com/Arman92/go-tdlib"
	log "github.com/sirupsen/logrus"
	"os"
	"telegram-parser/app"
	"telegram-parser/db"
)

func main() {
	checkENV()

	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")

	telegram := app.NewTgClient()
	telegram.Authorization()

	telegram.RunHandlingUpdates()

	cli, _ := db.ConnectToPostgres()
	//cli.GetAllData()
	cli.GetMessageById(-1001129804770, 408447614976)

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

func Log() {
	filename := "logfile.log"
	// Create the log file if doesn't exist. And append to it if it already exists.
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05" // "Mon Jan 2 15:04:05 MST 2006"
	Formatter.FullTimestamp = true

	log.SetFormatter(Formatter)
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	} else {
		log.SetOutput(f)
	}

	//log.Info("Some info. Earth is not flat")
	//log.Warning("This is a warning")
	//log.Error("Not fatal. An error. Won't stop execution")
	//log.Fatal("MAYDAY MAYDAY MAYDAY")
	//log.Panic("Do not panic")
}

package helpers

import (
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
	"os"
)

func CheckENV() {
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

func ConfigureLogrus() {
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

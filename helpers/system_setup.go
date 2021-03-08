package helpers

import (
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
	"os"
)

type env struct {
	name    string
	comment string
}

var envs = []env{
	{"PORT", "The port that the parser will listen to."},
	{"RABBITPORT", "Port on which rabbitMQ is listened."},
	{"PGHOST", "Host on which postgreSQL is listened."}, // "localhost"
	{"PGPORT", "Port on which postgreSQL is listened."}, // "5432"
	{"PGUSER", "PostgreSQL user name."},                 // "postgres"
	{"PGPASSWORD", "Password to access postgreSQL."},    // ***
	//{"PGDBNAME", "PostgreSQL database name."},           // "portgres" DEPRECATE
	{"TGTELNUMBER", "Phone number required to connect to the telegram client."},
	{"TGAPIID", "Application identifier for Telegram API access, which can be obtained at https://my.telegram.org   --- must be non-empty.."},
	{"TGAPIHASH", "Application identifier hash for Telegram API access, which can be obtained at https://my.telegram.org  --- must be non-empty.."},
}

// CheckEnv checks for all required global variables.
func CheckEnv() {
	for _, env := range envs {
		e := os.Getenv(env.name)
		if e == "" {
			logrus.Fatalf("Missing global variable %v. Usage: %v\n ", env.name, env.comment)
		}
	}
}

// ConfigureLogrus minimally configures logrus.
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

package flags

import (
	"flag"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Postgres *Postgres
	Telegram *Telegram
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	DbName   string
	Password string
}

type Telegram struct {
	TelephoneNumber string
	APIID           string
	APIHash         string
}

// Parse parses the command-line flags.
func Parse() *Config {
	var (
		postgresHost     string
		postgresPort     string
		postgresUser     string
		postgresDbName   string
		postgresPassword string

		telephoneNumber string
		apiID           string
		apiHash         string
	)

	flag.StringVar(&postgresHost, "pghost", "localhost", "host for connecting to postregres")
	flag.StringVar(&postgresPort, "pgport", "5432", "port for connecting to postregres")
	flag.StringVar(&postgresUser, "pguser", "postgres", "postregres user")
	flag.StringVar(&postgresDbName, "dbname", "telegram-parser", "postregres database name")
	flag.StringVar(&postgresPassword, "pwd", "", "password for postgres")

	flag.StringVar(&telephoneNumber, "tel", "", "phone number required to connect to the telegram client")
	flag.StringVar(&apiID, "id", "", "application identifier for Telegram API access, which can be obtained at https://my.telegram.org   --- must be non-empty..")
	flag.StringVar(&apiHash, "hash", "", "application identifier hash for Telegram API access, which can be obtained at https://my.telegram.org  --- must be non-empty..")
	flag.Parse()

	conf := &Config{
		&Postgres{
			Host:     postgresHost,
			Port:     postgresPort,
			User:     postgresUser,
			DbName:   postgresDbName,
			Password: postgresPassword,
		},
		&Telegram{
			TelephoneNumber: telephoneNumber,
			APIID:           apiID,
			APIHash:         apiHash,
		},
	}

	conf.checkValidity()

	return conf
}

//checkValidity checks validity of the fields
func (c *Config) checkValidity() {
	if c.Postgres.Password == "" {
		logrus.Panic("no password provided for postgres database")
	}

	if c.Telegram.TelephoneNumber == "" {
		logrus.Panic("no phone number provided to access telegrams to the client")
	}

	if c.Telegram.APIID == "" {
		logrus.Panic("no API ID provided to access telegrams to the client")
	}

	if c.Telegram.APIHash == "" {
		logrus.Panic("no API hash provided to access telegrams to the client")
	}
}

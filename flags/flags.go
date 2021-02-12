package flags

import (
	"flag"
	"github.com/sirupsen/logrus"
)

//if os.Getenv("POSTGRESPASSWORD") == "" {
//logrus.Panic(`Environment variable "POSTGRESPASSWORD" is blank`)
//}
//
//if os.Getenv("TGTELEPHONENUMBER") == "" {
//logrus.Panic(`Environment variable "TGTELEPHONENUMBER" is blank`)
//}
//
//if os.Getenv("TGAPIID") == "" {
//logrus.Panic(`Environment variable "TGAPIID" is blank`)
//}
//
//if os.Getenv("TGAPIHASH") == "" {
//logrus.Panic(`Environment variable "TGAPIHASH" is blank`)
//}

//("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", os.Getenv("POSTGRESPASSWORD"), "test")
//(`CREATE TABLE IF NOT EXISTS tg_parser ( message_id bigint, chat_id bigint, chat_title text, content text, date bigint, views integer, forwards integer, replies integer, PRIMARY KEY(message_id, chat_id) );`)

/*
host="localhost"    def
port=5432           def
user="postgres"  def
dbname="telegram-parser"   def
password=******
*/

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresDbName   string
	PostgresPassword string
}

// Parse
func Parse() *Config {
	var (
		postgresHost     string
		postgresPort     string
		postgresUser     string
		postgresDbName   string
		postgresPassword string
	)

	flag.StringVar(&postgresHost, "pghost", "localhost", "host for connecting to postregres") //-n=popopo
	flag.StringVar(&postgresPort, "pgport", "5432", "port for connecting to postregres")
	flag.StringVar(&postgresUser, "pguser", "postgres", "postregres user")
	flag.StringVar(&postgresDbName, "dbname", "telegram-parser", "postregres database name")
	flag.StringVar(&postgresPassword, "pwd", "", "password for postgres")
	flag.Parse()

	if postgresPassword == "" {
		logrus.Panic("no password provided for postgres database")
	}

	return &Config{
		PostgresHost:     postgresHost,
		PostgresPort:     postgresPort,
		PostgresUser:     postgresUser,
		PostgresDbName:   postgresDbName,
		PostgresPassword: postgresPassword,
	}
}

// ----------DEPRECATED----------
//// ParseFlags parses args and checks if the values is an ipv4 address
//func ParseFlags(addresses chan string) {
//	flag.Parse()
//	nodes := flag.Args()
//
//	for _, node := range nodes {
//		valid := isValidNodeAddress(node)
//		if !valid {
//			logrus.Errorf("The value `%v` passed in the arguments is not an address\n", node)
//			continue
//		}
//
//		addresses <- node
//	}
//}
//
//// isValidNodeAddress checks if value is an ipv4 address
//func isValidNodeAddress(hostport string) bool {
//	spl := strings.Split(hostport, ":")
//	if len(spl) != 2 {
//		return false
//	}
//
//	port := spl[1]
//	intPort, err := strconv.Atoi(port)
//	if err != nil {
//		return false
//	}
//
//	if intPort > 65535 || intPort < 1 {
//		return false
//	}
//
//	if spl[0] == "localhost" {
//		return true
//	}
//
//	spl = strings.Split(spl[0], ".")
//	if len(spl) != 4 {
//		return false
//	}
//
//	for _, host := range spl {
//		intHost, err := strconv.Atoi(host)
//		if err != nil {
//			return false
//		}
//
//		if intHost > 255 || intHost < 0 {
//			return false
//		}
//	}
//
//	return true
//}

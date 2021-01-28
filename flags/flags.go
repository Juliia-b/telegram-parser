package flags

import (
	"flag"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type Config struct {
	ParserAddrs []string // Addresses of parser services
}

func ParseFlags() *Config {
	config := &Config{}

	flag.Parse()
	nodes := flag.Args()

	for _, node := range nodes {
		valid := isValidNodeAddress(node)
		if !valid {
			logrus.Errorf("Service address %v is not valid\n", node)
			continue
		}

		config.ParserAddrs = append(config.ParserAddrs, node)
	}

	logrus.Infof("Service addresses: %#v\n", config.ParserAddrs)

	return config
}

func isValidNodeAddress(hostport string) bool {
	spl := strings.Split(hostport, ":")
	if len(spl) != 2 {
		return false
	}

	port := spl[1]
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	if intPort > 65535 || intPort < 1 {
		return false
	}

	if spl[0] == "localhost" {
		return true
	}

	spl = strings.Split(spl[0], ".")
	if len(spl) != 4 {
		return false
	}

	for _, host := range spl {
		intHost, err := strconv.Atoi(host)
		if err != nil {
			return false
		}

		if intHost > 255 || intHost < 0 {
			return false
		}
	}

	return true
}

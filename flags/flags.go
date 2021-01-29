package flags

import (
	"flag"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// ParseFlags parses args and checks if the values is an ipv4 address
func ParseFlags(addresses chan string) {
	flag.Parse()
	nodes := flag.Args()

	for _, node := range nodes {
		valid := isValidNodeAddress(node)
		if !valid {
			logrus.Errorf("The value `%v` passed in the arguments is not an address\n", node)
			continue
		}

		addresses <- node
	}
}

// isValidNodeAddress checks if value is an ipv4 address
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

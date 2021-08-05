package main

import (
	"fmt"
	"path"

	log "github.com/sirupsen/logrus"
)

type AuthCommand struct {
	Debug  bool   `short:"d" long:"debug" description:"Turn on debug output"`
	Config string `short:"c" long:"config" description:"Location of the configuration file" default:"/etc/logtailer/config.yaml"`
}

var authCommand AuthCommand

func (x *AuthCommand) Execute(args []string) error {

	if authCommand.Debug {
		setLogging(log.DebugLevel)
		log.Debugln("Setting log level to Debug")
	} else {
		setLogging(log.WarnLevel)
		log.Warn("Setting log level to Warn")
	}

	var config Config
	config.ReadConfig(x.Config)

	clientID := "cilogon:/client_id/3ea720f288b2762a8acf5f931eced0a6"
	scopes := []string{"openid"}

	code, err := RequestCode("https://cilogon-device-flow-proxy.herokuapp.com/device/code", clientID, scopes)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Copy code: %s\n", code.UserCode)
	fmt.Printf("then open: %s\n", code.VerificationURI)
	fmt.Printf("Device Code: %s\n", code.DeviceCode)

	accessToken, err := PollToken("https://cilogon-device-flow-proxy.herokuapp.com/device/token", clientID, code)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Access token: %s\n", accessToken.AccessToken)

	// Write the Access Token and Refresh Token
	path.Join(config)

	return nil

}

func init() {
	parser.AddCommand("auth",
		"Aqcuire authentication credentials",
		"The auth command is used to acquire authentication credentials",
		&authCommand)
}

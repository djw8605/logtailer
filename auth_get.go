package main

import (
	"fmt"
	"net/http"

	device "github.com/cli/oauth/device"
	log "github.com/sirupsen/logrus"
)

type AuthCommand struct {
	Debug bool `short:"d" long:"debug" description:"Turn on debug output"`
}

var authCommand AuthCommand

func (x *AuthCommand) Execute(args []string) error {

	if authCommand.Debug {
		setLogging(log.DebugLevel)
	} else {
		setLogging(log.WarnLevel)
	}

	fmt.Printf("Adding (all=): %#v\n", args)
	clientID := "OSG-TOKEN-LOGTAILER"
	scopes := []string{"my_rabbit_server.write:osg-htcondor-xfer/osg-htcondor-xfer"}
	httpClient := http.DefaultClient

	code, err := device.RequestCode(httpClient, "https://cilogon-device-flow-proxy.herokuapp.com/device/code", clientID, scopes)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Copy code: %s\n", code.UserCode)
	fmt.Printf("then open: %s\n", code.VerificationURI)

	accessToken, err := device.PollToken(httpClient, "https://cilogon-device-flow-proxy.herokuapp.com/device/token", clientID, code)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Access token: %s\n", accessToken.Token)

	return nil
}

func init() {
	parser.AddCommand("auth",
		"Aqcuire authentication credentials",
		"The auth command is used to acquire authentication credentials",
		&authCommand)
}

package main

import (
	log "github.com/sirupsen/logrus"
)

// Start the AMQP connection
func StartAmqp(config Config, initialCred string, msgChannel <-chan string, credChannel <-chan string) {

	// Initialize the connection
	for true {
		select {
		case msg := <-msgChannel:
			log.Debugln("Received message:", msg)
		case cred := <-credChannel:
			// Close the current AMQP connection and create a another with the credentail
			log.Debugln("New Credential:", cred)

		}
	}

}

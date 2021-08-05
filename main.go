package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	gocron "github.com/go-co-op/gocron"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

//type WatchFiles struct {
//	Files string `positional-arg-name:"file"`
//}

type Options struct {
	// Turn on the debug logging
	Debug bool `short:"d" long:"debug" description:"Turn on debug logging"`

	// Specify the configuration file
	Config string `short:"c" long:"config" description:"Location of the configuration file" default:"/etc/logtailer/config.yaml"`

	// Positional arguemnts
	//Files WatchFiles `description:"Files to watch" positional-args:"1"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

// main
func main() {

	var err error
	parser.SubcommandsOptional = true
	if _, err = parser.Parse(); err != nil {
		// err will tell us if the user just did a "-h"
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	if options.Debug {
		setLogging(log.DebugLevel)
		log.Debug("Setting loglevel to Debug")
	} else {
		setLogging(log.WarnLevel)
	}

	// Read in the configuration file
	var config Config
	config.ReadConfig(options.Config)
	log.Debugf(config.Amqp.Host)
	s := gocron.NewScheduler(time.UTC)
	messageBus := make(chan string)
	credentialBus := make(chan string)
	initialCred := "Blah"
	go StartAmqp(config, initialCred, messageBus, credentialBus)
	s.Every(1).Second().SingletonMode().Do(renewCred, credentialBus)
	s.StartAsync()

	// Read in the watched files from the yaml
	watchedFile := "./stuff"
	//watchedFile := args[0]
	log.Debugf("Watched file: %s", watchedFile)

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add("./stuff"); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}

func setLogging(logLevel log.Level) error {
	textFormatter := log.TextFormatter{}
	textFormatter.DisableLevelTruncation = true
	textFormatter.FullTimestamp = true
	log.SetFormatter(&textFormatter)
	log.SetLevel(logLevel)
	return nil
}

func renewCred(param chan<- string) {
	param <- "Hello"
	log.Debugln("Renewing cred")
}

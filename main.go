package main

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	auth 
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
	auth.init()
	var err error
	if _, err = parser.Parse(); err != nil {
		// err will tell us if the user just did a "-h"
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if options.Debug {
		setLogging(log.DebugLevel)
	} else {
		setLogging(log.WarnLevel)
	}

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
	return nil
}

package main

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Amqp struct {
		User string
		Host string
	}
	CredDir string `yaml:creddir`
}

// Read in the configuration file
func (c *Config) ReadConfig(configFile string) (*Config, error) {
	// Read the file
	log.Debugf("Reading in the file: %s", configFile)
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
		return c, err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return c, err
	}

	return c, nil

}

package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ClientID       string `yaml:"client_id"`
	ClientSecret   string `yaml:"client_secret"`
	PlaylistTarget string `yaml:"playlist"`
	Port           string `yaml:"port"`
}

func (c *Config) parseConfig(path string) {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(`Error reading config, make sure you have config file 
		named musync.yaml or point 
		to another config using the -c flag`)
	}
	yaml.Unmarshal([]byte(fileData), &c)

}

package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Conf the representation of a config file
type Conf struct {
	// Token the bot token
	Token string `json:"token"`
	// Prefix the prefix used to issue commands to the bot
	Prefix string `json:"prefix"`
}

// Config the global config variable
var Config *Conf

// NewConfig loads the config from the config file
func NewConfig(p string) *Conf {
	// Make sure path exists
	if _, err := os.Stat(p); err != nil {
		Die(err)
	}
	// Load config from file
	f, _ := ioutil.ReadFile(p)
	var cfg = &Conf{}
	err := yaml.Unmarshal(f, &cfg)
	Die(err)

	return cfg
}

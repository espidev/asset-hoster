package main

import (
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	Port int
	Secret string
	Domain string
	Debug bool
}

const (
	ConfigLocation = "./config.toml"
	DefaultConfig  = `
Port = 3000
Secret = "notsecret"
Domain = "localhost"
Debug = true
`)

func setupConfig() {
	if _, err := os.Stat(ConfigLocation); os.IsNotExist(err) {
		err := ioutil.WriteFile(ConfigLocation, []byte(DefaultConfig), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	if _, err := toml.DecodeFile(ConfigLocation, &config); err != nil {
		log.Fatal(err)
	}
}
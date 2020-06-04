package main

import (
	"log"

	"github.com/jinzhu/configor"
)

//Config contains global app's configuration
var Config AppConfig

type LoginSource struct {
	ID            int
	Name          string
	Strategy      string
	OptionalEmail string
	URL           string
	Key1          string
	Key2          string
}

// AppConfig contains app's configuration
type AppConfig struct {
	Server struct {
		Public string `default:"http://localhost:8040"`
		Port   string `default:":80"`
		Data   string `default:"./data"`
	}
	DB struct {
		Path     string
		User     string
		Host     string
		Password string
		Database string
	}
	Login []LoginSource
}

//LoadFromFile method loads and parses config file
func (c *AppConfig) LoadFromFile(url string) {
	err := configor.Load(&Config, url)
	if err != nil {
		log.Fatalf("Can't load the config file: %s", err)
	}
}

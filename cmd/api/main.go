package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/standardWebServer/internal/app/api"
)

var (
	configPath string = "configs/api.toml"
)

func init() {
	flag.StringVar(&configPath, "path", "configs/api.toml", "path to config file in .toml format")
}

func main() {
	flag.Parse()
	log.Println("It works")
	config := api.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Println("can not find configs file. using default values:", err)
	}

	server := api.New(config)

	log.Fatal(server.Start())
}

package main

import (
	"log"

	"github.com/chiaf1/plclogger/internal/config"
)

func main() {
	var conf config.Config
	err := conf.Load("./config.yaml")
	if err != nil {
		panic(err)
	}
	log.Println("Config loaded")
}

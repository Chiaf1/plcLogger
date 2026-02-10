package main

import (
	"fmt"
	"log"

	"github.com/chiaf1/plclogger/internal/config"
	datastorage "github.com/chiaf1/plclogger/internal/dataStorage"
)

func main() {
	var conf config.Config
	err := conf.Load("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config loaded")

	var lv datastorage.LastValues
	err = lv.LoadLastValues("./data/last_values.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Last values loaded")

	fmt.Println(lv)

}

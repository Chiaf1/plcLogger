package main

import (
	"log"

	"github.com/chiaf1/plclogger/internal/config"
	datastorage "github.com/chiaf1/plclogger/internal/dataStorage"
)

func main() {
	//loading config
	var conf config.Config
	err := conf.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config loaded")

	//loading last values
	var lv datastorage.LastValues
	err = lv.LoadLastValues("./data/last_values.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Last values loaded")

	//retreving the current values
	// dataLogOnChange := conf.GetOnChangeLogTags()

}

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
	lv.SetValue("tag1", 123)
	lv.SetValue("tag2", 123.456)
	lv.SetValue("tag3", "Luigi")
	lv.SetValue("tag4", true)
	lv.SetValue("tag5", nil)
	fmt.Println(lv)
	lv.SaveLastValues("./data/last_values.json")

}

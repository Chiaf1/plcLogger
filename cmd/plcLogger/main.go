package main

import (
	"log"

	"github.com/chiaf1/plclogger/internal/config"
	datastorage "github.com/chiaf1/plclogger/internal/dataStorage"
	"github.com/chiaf1/plclogger/internal/logger"
	plccomunication "github.com/chiaf1/plclogger/internal/plcComunication"
)

const CONFIG_PATH = "./config.yaml"
const LAST_VALUES_PATH = "./data/last_values.json"
const ON_CHANGE_LOG_PATH = "./log/onChange.log"
const PERIODIC_LOG_PATH = "./log/periodic.log"

func main() {
	//loading config
	var conf config.Config
	err := conf.Load(CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config loaded")

	//loading last values
	var lv datastorage.LastValues
	err = lv.LoadLastValues(LAST_VALUES_PATH)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Last values loaded")

	//retreving the current values and checking if they changed
	var rv plccomunication.ReadValues
	err = rv.UpdateCurrentVals(conf.DataToLog, conf.Connection)
	if err != nil {
		log.Fatal(err)
	}
	logger.CheckChangedValues(lv, rv.GetOnChange(), LAST_VALUES_PATH, ON_CHANGE_LOG_PATH)
	log.Println("OnChange logged")

	logger.LogPeriodic(rv.GetPeriodic(), PERIODIC_LOG_PATH)
	log.Println("Periodic logged")

}

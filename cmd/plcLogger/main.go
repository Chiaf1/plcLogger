package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

const CONFIG_PATH = "./config.yaml"
const LAST_VALUES_PATH = "./data/last_values.json"
const ON_CHANGE_LOG_PATH = "./log/onChange.log"
const PERIODIC_LOG_PATH = "./log/periodic.log"

func setupSignalHandler(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()
}

func main() {
	// Creating the parent context for all the routines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupSignalHandler(cancel)

	<-ctx.Done()
	/*
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
	*/
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chiaf1/plclogger/internal/config"
	datastorage "github.com/chiaf1/plclogger/internal/dataStorage"
	"github.com/chiaf1/plclogger/internal/logger"
	plccomunication "github.com/chiaf1/plclogger/internal/plcComunication"
	plcdrivers "github.com/chiaf1/plclogger/internal/plcDrivers"
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

	// Loading config from file or creating one if not present
	var conf config.Config
	err := conf.Load(CONFIG_PATH)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	log.Println("config loaded successfully")

	// Load last values
	var lv datastorage.LastValues
	err = lv.LoadLastValues(LAST_VALUES_PATH)
	if err != nil {
		log.Fatalf("error loading last values: %v", err)
	}
	log.Println("Last values loaded successfolly")

	// Creating new plc driver
	plc, err := plcdrivers.NewPLCDriver(conf.Connection)
	if err != nil {
		log.Fatalf("error creating the PLC comunication driver: %v", err)
	}

	// Creating ReadValues
	rv := plccomunication.NewReadValues()
	rv.ReplaceAllTags(conf.DataToLog) // Load all tags

	// Start PLC comunication
	plc.Connect()
	defer plc.Disconnect()

	// Start plc polling
	rv.StartPoller(ctx, plc, conf.App.OnChangeClock, log.Printf)

	// Start web server
	if conf.App.EnableWebServer {
		// at the moment not present yet
		// go startWebServer(ctx, rv)
	}

	time.Sleep(5 * time.Second) // sleep 5s so that the poller has time to updated the values before the logger start to read them

	// Start periodic logger
	go func() {
		// start ticker for periodic interval
		ticker := time.NewTicker(conf.App.PeriodicLogInterval)
		defer ticker.Stop()

		// Logging loop
		for {
			snap := rv.SnapShot()
			cv := make(logger.CurrentValues)
			for _, t := range snap {
				if t.PeriodicLog {
					cv[t.Name] = t.Val
				}
			}
			err := logger.LogPeriodic(cv, PERIODIC_LOG_PATH)
			if err != nil {
				log.Printf("periodic log error: %v", err)
			}

			// Check log file rotatio and archive
			err = logger.CheckArchiveRotation(
				PERIODIC_LOG_PATH,
				conf.App.PeriodicLog.MaxSize,
				conf.App.PeriodicLog.MaxAge,
				conf.App.PeriodicLog.ArchivePath,
				conf.App.PeriodicLog.ArchiveMaxSize,
			)
			if err != nil {
				log.Printf("erro rotating periodic log file: %v", err)
			}

			select {
			case <-ctx.Done():
				// shutdown
				return
			case <-ticker.C:
				// Restart loop
				continue
			}
		}
	}()

	// Start onChange logger
	if conf.App.EnableOnChangeLog {
		go func() {
			// Start ticker
			ticker := time.NewTicker(conf.App.OnChangeClock)
			defer ticker.Stop()

			// Logging loop
			for {
				snap := rv.SnapShot()
				cv := make(logger.CurrentValues)
				for _, t := range snap {
					if t.OnChangeLog {
						cv[t.Name] = t.Val
					}
				}
				err := logger.CheckChangedValues(lv, cv, LAST_VALUES_PATH, ON_CHANGE_LOG_PATH)
				if err != nil {
					log.Printf("onChange log error: %v", err)
				}

				// Check log file rotatio and archive
				err = logger.CheckArchiveRotation(
					ON_CHANGE_LOG_PATH,
					conf.App.OnChangeLog.MaxSize,
					conf.App.OnChangeLog.MaxAge,
					conf.App.OnChangeLog.ArchivePath,
					conf.App.OnChangeLog.ArchiveMaxSize,
				)
				if err != nil {
					log.Printf("erro rotating OnChange log file: %v", err)
				}

				select {
				case <-ctx.Done():
					//Shutdown
					return
				case <-ticker.C:
					// Restart loop
					continue
				}
			}
		}()
	}

	// Wait for stop singal from system
	<-ctx.Done()

}

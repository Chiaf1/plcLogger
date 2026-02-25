package plccomunication

import (
	"context"
	"time"

	"github.com/chiaf1/plclogger/internal/domain"
)

// StartPolling starts a loop that periodically updates the value of the tags registered in ReadValues.
// - ctx: stops the polling cleanly
// - plc: drvier PLC
// - interval: polling interval
// - logf: logging function (es.log.Printf)
func (rv *ReadValues) StartPolling(ctx context.Context, plc domain.PLCDriver, interval time.Duration, logf func(string, ...any)) {
	// polling interval tick creation
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// first tag updated at the start
	if err := rv.UpdateCurrentVals(plc); err != nil {
		logf("plc initial read warning: %v", err)
	}

	// polling loop
	for {
		select {
		case <-ctx.Done():
			// shutdown
			return
		case <-ticker.C:
			// periodic tag update
			if err := rv.UpdateCurrentVals(plc); err != nil && logf != nil {
				// the loop doesn't stop but we log the error
				logf("plc read waring: %v", err)
			}
		}
	}
}

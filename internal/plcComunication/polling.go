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

// StartPoller creates an internal context adn starts StartPolling in a go routine
// it protects over starting multiple pollers, if the poller already started it returns false
func (rv *ReadValues) StartPoller(plc domain.PLCDriver, interval time.Duration, logf func(string, ...any)) bool {
	rv.pollerMu.Lock()
	defer rv.pollerMu.Unlock()

	if rv.polling {
		return false
	}

	// shutdown context StopPoller or parent ctx
	ctx, cancel := context.WithCancel(context.Background())
	rv.cancelPoll = cancel
	rv.polling = true

	go func() {
		defer func() {
			rv.pollerMu.Lock()
			rv.polling = false
			rv.cancelPoll = nil
			rv.pollerMu.Unlock()
		}()
		rv.StartPolling(ctx, plc, interval, logf)
	}()

	return true
}

// StopPoller stop polling if active. Returns true if stopped, false if already stopped
func (rv *ReadValues) StopPoller() bool {
	rv.pollerMu.Lock()
	defer rv.pollerMu.Unlock()

	if !rv.polling || rv.cancelPoll == nil {
		return false
	}
	rv.cancelPoll()
	return true
}

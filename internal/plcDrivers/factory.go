package plcdrivers

import (
	"fmt"

	"github.com/chiaf1/plclogger/internal/domain"
)

// NewPLCDriver returns a plc driver based on the config passed, if the driver in the config is not supported an error is returned
func NewPLCDriver(conf domain.ConnectionConfig) (domain.PLCDriver, error) {
	switch conf.Protocol {
	case "s7":
		return NewS7Client(conf), nil
	case "opcua":
		return NewOPCUAClient(conf), nil
	default:
		return nil, fmt.Errorf("PLC comunication protocol %q, not supported", conf.Protocol)
	}
}

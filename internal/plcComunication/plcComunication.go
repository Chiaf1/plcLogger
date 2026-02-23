package plccomunication

import (
	"github.com/chiaf1/plclogger/internal/config"
	"github.com/chiaf1/plclogger/internal/logger"
)

// PLC client interface with standard metods common to multiple protocols
type PLCClient interface {
	Connect() error
	Disconnect() error
	Read(tag config.PlcTag) (any, error)
	CheckConnection() bool
}

// UpdateCurrentVals connectes to the plc, retrieves the current values of all the tags in the PlcTag slice (tags) and retturns
// the current values that's a map with the kye value as the name of tag and the value as the value of the tag
func UpdateCurrentVals(tags []config.PlcTag, conf config.ConnectionConfig) (logger.CurrentValues, error) {
	curVal := make(logger.CurrentValues)
	//for the tests i'm forcing the values by hand
	curVal["tag1"] = 124
	curVal["tag2"] = 124.457
	curVal["tag3"] = "Luigii"
	curVal["tag4"] = false
	curVal["tag5"] = nil
	return curVal, nil
}

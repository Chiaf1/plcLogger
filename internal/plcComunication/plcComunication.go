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

// ReadValues type is a slice of Value that rapresent a tag read from the PLC with flags related to the type of log they belong to
type ReadValues []Value

// Value type that raperesent one tag with the flags PeriodicLog, OnChangeLog, ShowDashboard to rapresent the log they belong to
type Value struct {
	Name          string
	Val           any
	PeriodicLog   bool
	OnChangeLog   bool
	ShowDashboard bool
}

// UpdateCurrentVals connectes to the plc, retrieves the current values of all the tags of the slice tags and stores them in the struct
// ReadValues. It read all the values from the PLC then you can filter them based on what you need
func (rv *ReadValues) UpdateCurrentVals(tags []config.Tag, conf config.ConnectionConfig) error

// GetPeriodic from the ReadVals slice retrive only the vals to log periodically
func (rv ReadValues) GetPeriodic() logger.CurrentValues {
	LogTags := make(logger.CurrentValues)
	for _, v := range rv {
		if v.PeriodicLog {
			LogTags[v.Name] = v.Val
		}
	}
	return LogTags
}

// GetOnChange from the ReadVals slice retrive only the vals to log on value change
func (rv ReadValues) GetOnChange() logger.CurrentValues {
	LogTags := make(logger.CurrentValues)
	for _, v := range rv {
		if v.OnChangeLog {
			LogTags[v.Name] = v.Val
		}
	}
	return LogTags
}

// GetShowDashboard from the ReadVals slice retrive only the vals to show on the dashboard
func (rv ReadValues) GetShowDashboard() logger.CurrentValues {
	LogTags := make(logger.CurrentValues)
	for _, v := range rv {
		if v.ShowDashboard {
			LogTags[v.Name] = v.Val
		}
	}
	return LogTags
}

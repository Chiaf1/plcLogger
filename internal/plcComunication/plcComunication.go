package plccomunication

import (
	"fmt"

	"github.com/chiaf1/plclogger/internal/domain"
	"github.com/chiaf1/plclogger/internal/logger"
	plcdrivers "github.com/chiaf1/plclogger/internal/plcDrivers"
)

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
func (rv *ReadValues) UpdateCurrentVals(tags []domain.Tag, conf domain.ConnectionConfig) error {
	if len(tags) <= 0 {
		return fmt.Errorf("error no tags to read")
	}

	client := plcdrivers.NewS7Client(conf)

	err := client.Connect()
	if err != nil {
		return fmt.Errorf("error connecting to plc client: %w", err)
	}
	defer client.Disconnect()

	var errs []error
	for _, t := range tags {
		val, err := client.Read(t.PlcTag)
		if err != nil {
			errs = append(errs, fmt.Errorf("error reading tag %s from plc: %w", t.Name, err))
			continue
		}
		newVal := Value{
			Name:          t.Name,
			Val:           val,
			PeriodicLog:   t.PeriodicLog,
			OnChangeLog:   t.OnChangeLog,
			ShowDashboard: t.ShowDashboard,
		}
		err = rv.UpdateVal(newVal)
		if err != nil {
			errs = append(errs, fmt.Errorf("error updating tag value: %w", err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("read completed with %d errors: %v", len(errs), errs)
	}
	return nil
}

// UpdateVal updates the value if it exist, if not it appends it
func (rv *ReadValues) UpdateVal(v Value) error {
	if v.Name == "" {
		return fmt.Errorf("tag with no name can't be added/updated")
	}
	if len(*rv) > 0 {
		for i, p := range *rv {
			if p.Name == v.Name {
				(*rv)[i] = v
				return nil
			}
		}
	}
	*rv = append(*rv, v)
	return nil
}

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

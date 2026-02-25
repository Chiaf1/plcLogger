package plccomunication

import (
	"fmt"
	"maps"
	"sync"

	"github.com/chiaf1/plclogger/internal/domain"
)

// ReadValues is a struct containing a map to all the Values of the tags that are read from the PLC
// this struct also has a mutex to access the data safely from different threads
// the value map has the key value queal to the name of the tag
type ReadValues struct {
	mu   sync.RWMutex
	vals map[string]Value
}

// Value type that raperesent one tag with the flags PeriodicLog, OnChangeLog, ShowDashboard to rapresent the log they belong to
type Value struct {
	Name          string
	Val           any
	PeriodicLog   bool
	OnChangeLog   bool
	ShowDashboard bool
}

// NewReadValues is a constructor for ReadValues struct
func NewReadValues() *ReadValues {
	return &ReadValues{vals: make(map[string]Value)}
}

// UpdateVal is used to safely update the value of a tag or add it if not present
func (rv *ReadValues) UpdateVal(v Value) error {
	if v.Name == "" {
		return fmt.Errorf("tag with no name can't be added/updated")
	}
	rv.mu.Lock()
	rv.vals[v.Name] = v
	rv.mu.Unlock()
	return nil
}

// UpdateCurrentVals connectes to the plc, retrieves the current values of all the tags of the slice tags and stores them in the struct
// ReadValues. It updates the structure ReadValues in thread safe way
func (rv *ReadValues) UpdateCurrentVals(tags []domain.Tag, plc domain.PLCDriver) error {
	if len(tags) <= 0 {
		return fmt.Errorf("error no tags to read")
	}

	// We'll work on a temp structure so that we can lock the real one only after we have read all the tags
	tmp := make(map[string]Value, len(tags))
	// we create a slice to store all the errors during the reading of the tags
	var errs []error
	for _, t := range tags {
		val, err := plc.Read(t.PlcTag)
		if err != nil {
			errs = append(errs, fmt.Errorf("error reading tag %s from plc: %w", t.Name, err))
			continue
		}
		tmp[t.Name] = Value{
			Name:          t.Name,
			Val:           val,
			PeriodicLog:   t.PeriodicLog,
			OnChangeLog:   t.OnChangeLog,
			ShowDashboard: t.ShowDashboard,
		}
	}

	// now we'll update the real structure
	rv.mu.Lock()
	maps.Copy(rv.vals, tmp)
	rv.mu.Unlock()

	if len(errs) > 0 {
		return fmt.Errorf("read completed with %d errors: %v", len(errs), errs)
	}
	return nil
}

// Snapshot safely retrives a copy/snapshot of the current state of the values stored in the ReadValue structure
func (rv *ReadValues) SnapShot() []Value {
	rv.mu.RLock()
	defer rv.mu.RUnlock()

	out := make([]Value, 0, len(rv.vals))
	for _, v := range rv.vals {
		out = append(out, v)
	}
	return out
}

// GetByName safely retrives the value of a single tag based on the name
func (rv *ReadValues) GetByName(name string) (Value, bool) {
	rv.mu.RLock()
	v, ok := rv.vals[name]
	rv.mu.RUnlock()
	return v, ok
}

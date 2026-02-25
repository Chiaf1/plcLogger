package plccomunication

import (
	"fmt"
	"log"
	"maps"
	"sync"

	"github.com/chiaf1/plclogger/internal/domain"
)

// ReadValues is a struct containing a map to all the Values of the tags that are read from the PLC
// this struct also has a mutex to access the data safely from different threads
// the value map has the key value queal to the name of the tag
type ReadValues struct {
	mu   sync.RWMutex
	tags map[string]*Value
}

// Value type that raperesent one tag with the flags PeriodicLog, OnChangeLog, ShowDashboard to rapresent the log they belong to
type Value struct {
	domain.Tag
	Val any
}

// NewReadValues is a constructor for ReadValues struct
func NewReadValues() *ReadValues {
	return &ReadValues{tags: make(map[string]*Value)}
}

// AddOrUpdateVal is used to safely update the value of a tag or add it if not present
func (rv *ReadValues) AddOrUpdateVal(v Value) error {
	if v.Name == "" {
		return fmt.Errorf("tag with no name can't be added/updated")
	}
	rv.mu.Lock()
	defer rv.mu.Unlock()

	if existing, ok := rv.tags[v.Name]; ok {
		existing.Tag = v.Tag
		// existing.val is keept untouched, other functions will update it
		return nil
	}
	rv.tags[v.Name] = &Value{
		Tag: v.Tag,
		Val: nil, //no read yet
	}
	return nil
}

// RemoveTag removes the tag from the registry based on the name
func (rv *ReadValues) RemoveTag(name string) {
	rv.mu.Lock()
	delete(rv.tags, name)
	rv.mu.Unlock()
}

// ReplaceAllTags replace the current tag registry with a new one based on the tag list passed as argument
func (rv *ReadValues) ReplaceAllTags(tags []domain.Tag) {
	rv.mu.Lock()
	defer rv.mu.Unlock()

	next := make(map[string]*Value, len(tags))
	for _, t := range tags {
		if t.Name == "" {
			continue
		}
		// if already present just updates it
		if existing, ok := rv.tags[t.Name]; ok {
			next[t.Name] = &Value{
				Tag: t,
				Val: existing.Val,
			}
			continue
		}
		next[t.Name] = &Value{Tag: t}
	}
	rv.tags = next
}

// UpdateCurrentVals connectes to the plc and updates all the values of the ReadValue tag registry
func (rv *ReadValues) UpdateCurrentVals(plc domain.PLCDriver) error {
	// first let's create a snapshot to work on it safely without locking all functions
	rv.mu.RLock()
	tmp := make(map[string]*Value, len(rv.tags))
	maps.Copy(tmp, rv.tags)
	rv.mu.RUnlock()

	// we create a slice to store all the errors during the reading of the tags
	var errs []error
	for name, entry := range tmp {
		val, err := plc.Read(entry.PlcTag)
		if err != nil {
			log.Printf("error reading tag %s from plc: %v", name, err)
			errs = append(errs, fmt.Errorf("error reading tag %s from plc: %w", name, err))
			continue
		}
		// updates the value safely
		rv.mu.Lock()
		if current, ok := rv.tags[name]; ok {
			current.Val = val
		}
		rv.mu.Unlock()
	}

	if len(errs) > 0 {
		return fmt.Errorf("read completed with %d errors: %v", len(errs), errs)
	}
	return nil
}

// Snapshot safely retrives a copy/snapshot of the current state of the values stored in the ReadValue structure
func (rv *ReadValues) SnapShot() []Value {
	rv.mu.RLock()
	defer rv.mu.RUnlock()

	out := make([]Value, 0, len(rv.tags))
	for _, v := range rv.tags {
		out = append(out, *v)
	}
	return out
}

// GetByName safely retrives the value of a single tag based on the name
func (rv *ReadValues) GetByName(name string) (Value, bool) {
	rv.mu.RLock()
	defer rv.mu.RUnlock()
	v, ok := rv.tags[name]
	if !ok {
		return Value{}, false
	}
	return *v, ok
}

// Names gets all the names of the tags in the registry
func (rv *ReadValues) Names() []string {
	rv.mu.RLock()
	defer rv.mu.RUnlock()
	out := make([]string, 0, len(rv.tags))
	for name := range rv.tags {
		out = append(out, name)
	}
	return out
}

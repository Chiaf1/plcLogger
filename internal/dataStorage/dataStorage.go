package datastorage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/chiaf1/plclogger/internal/utils"
)

type LastValues map[string]Vals

type Vals struct {
	Val any        `json:"value"`
	Ts  *time.Time `json:"timestamp"`
}

// LoadLastValues laods the last values from file
func (lv *LastValues) LoadLastValues(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("Error reading the last value file: %w", err)
	}
	err = json.Unmarshal(data, lv)
	if err != nil {
		return fmt.Errorf("Error during parsing of JSON last values file: %w", err)
	}
	return nil
}

// SaveLastValues saves the last values to the "path" using the WriteFileAtomic function
func (lv *LastValues) SaveLastValues(path string) error {
	data, err := json.Marshal(lv)
	if err != nil {
		return fmt.Errorf("Error while parsing to JSON: %w", err)
	}
	return utils.WriteFileAtomic(path, data, 0644)
}

// SetValue sets or adds a value to the LastValues map. if present it updates it if not it creates it
func (lv *LastValues) SetValue(name string, val any) error {
	if name == "" {
		return fmt.Errorf("Name must have a value")
	}
	// inizialize the map if needed
	if *lv == nil {
		*lv = make(LastValues)
	}
	ts := time.Now()

	v, exists := (*lv)[name]
	if exists {
		//update value
		v.Val = val
		v.Ts = &ts
		(*lv)[name] = v
		return nil
	}

	// create new
	(*lv)[name] = Vals{
		Val: val,
		Ts:  &ts,
	}
	return nil
}

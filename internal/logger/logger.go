package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	datastorage "github.com/chiaf1/plclogger/internal/dataStorage"
)

// CurrentValues is a struct rappresenting the current value of the corresponding tag (the key)
type CurrentValues map[string]any

// LogOnChangeType is the type used in the log file for the on change values. It has:
// Ts: time stamp, Tag: the name of the tag, Old: the old value, New: the new value
type LogOnChangeType struct {
	Ts  string `json:"ts"`
	Tag string `json:"tag"`
	Old any    `json:"old"`
	New any    `json:"new"`
}

// CheckChangedValues this function ranges over the current values and compares them to the LastValues, if the new ones are different
// or the LastValues don't have them stored it logs the difference and updates the LastValues. Once it has ranged over all the values
// it saves the LastValues to file
func CheckChangedValues(lv datastorage.LastValues, cv CurrentValues, pathLv string, pathLog string) {
	changed := false
	//if the list of old values is empty i create it
	if lv == nil {
		lv = make(datastorage.LastValues)
	}
	for k, curVal := range cv {
		last, exists := lv[k]
		//it checks if the key exist in the old values, if it doesn't it adds and logs it
		if !exists {
			LogOnChange(pathLog, k, nil, curVal)
			lv.SetValue(k, curVal)
			changed = true
			continue
		}
		if !SmartEqual(last.Val, curVal) {
			LogOnChange(pathLog, k, last.Val, curVal)
			lv.SetValue(k, curVal)
			changed = true
		}
	}
	// it saves the last value file if something changed
	if changed {
		lv.SaveLastValues(pathLv)
	}
}

// LogOnChange creates the json format of the data to log and passes it to AppendToLogFile to actually log it.
func LogOnChange(path string, name string, oldVal any, newVal any) error {
	// for the timestamp it uses RFC3339 to have a more readable time stamp
	ts := time.Now().Format(time.RFC3339)

	newLog := LogOnChangeType{
		Ts:  ts,
		Tag: name,
		Old: oldVal,
		New: newVal,
	}

	data, err := json.Marshal(newLog)
	if err != nil {
		return fmt.Errorf("Error parsing new values to JSON: %w", err)
	}

	return AppendToLogFile(path, append(data, '\n'))
}

// AppendToLogFile it appends the slice of bytes to the log file creating a new line. If the fiel or directory doens't
// exists it creates it
func AppendToLogFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	//create the directory, if it already exists nothing happens
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("Cannot create directory: %w", err)
	}
	// Now it opens the file with the append flag, writes the data to it and then closes it
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Error opening the file %s: %w", path, err)
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("Error writing file: %w", err)
	}
	return nil
}

// ToFloat64 thsi function normalizez numeric values in float64. This can be usefull when comparing a numeric value read from a json
// that are all normalized in float64 with a numeric value coming from a different part of the program that uses other types.
// the function returns the float value and a bool to show if it was converter (true) or not (false)
func ToFloat64(v any) (float64, bool) {
	switch x := v.(type) {
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(x).Int()), true
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(x).Uint()), true
	case float32, float64:
		return reflect.ValueOf(x).Float(), true
	default:
		return 0, false
	}
}

// SmartEqual this function compares two values, if they are numric they get both normalized to float64 and compared
// if they are not numeric they are compared with reflect.DeepEaqual
func SmartEqual(a, b any) bool {
	// normalize numeric values to float64
	af, okA := ToFloat64(a)
	bf, okB := ToFloat64(b)
	// if the values are numeric confronts the normalized values
	if okA && okB {
		// here a tolleranze can be added for float values
		// const eps = 1e-9
		// return math.Abs(af-bf) < eps
		return af == bf
	}

	// if the values are not numeric deep equal is called
	return reflect.DeepEqual(a, b)

}

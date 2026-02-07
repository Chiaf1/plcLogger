package logger

import (
	"reflect"

	"github.com/chiaf1/plclogger/internal/datastorage"
)

type CurrentValues map[string]any

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
			lv.AddValue(k, curVal)
			changed = true
			continue
		}
		if !reflect.DeepEqual(last.Val, curVal) {
			LogOnChange(pathLog, k, last.Val, curVal)
			lv.AddValue(k, curVal)
			changed = true
		}
	}
	// it saves the last value file if something changed
	if changed {
		lv.SaveLastValues(pathLv)
	}
}

func LogOnChange(path string, name string, oldVal any, newVal any)

package domain

/*
The package Domain is used to store all types' definitions that are common to multiple packages

This helps to avoid import loops inside the code
*/

import "time"

// ConnectionConfig defines the configurations used to connect to the PLC
type ConnectionConfig struct {
	Ip              string        `yaml:"ip"`
	Rack            int           `yaml:"rack"`
	Slot            int           `yaml:"slot"`
	Protocol        string        `yaml:"protocol"`
	Timeout         time.Duration `yaml:"timeout"`
	ConnectionRetry int           `yaml:"connectionRetry"`
}

// AppConfig stores the configurations of the app behavior
type AppConfig struct {
	PeriodicLogInterval time.Duration `yaml:"periodicLogInterval"`
	OnChangeClock       time.Duration `yaml:"onChangeClock"`
	EnableWebServer     bool          `yaml:"enableWebServer"`
	EnableOnChangeLog   bool          `yaml:"enableOnChangeLog"`
	PeriodicLog         LogConf       `yaml:"periodicLog"`
	OnChangeLog         LogConf       `yaml:"onChangeLog"`
}

// Tag defines all data related to one tag
type Tag struct {
	PlcTag        `yaml:",inline"`
	PeriodicLog   bool `yaml:"periodicLog"`
	OnChangeLog   bool `yaml:"onChangeLog"`
	ShowDashboard bool `yaml:"showDashboard"`
}

// PlcTag defines the tag data strictly related to the plc
type PlcTag struct {
	Name  string       `yaml:"name"`
	Type  PlcType      `yaml:"type"`
	S7    S7Mapping    `yaml:"s7"`
	OPCUA OPCUAMapping `yaml:"opca"`
}

// This enum represents the data types that can be read from the PLC
type PlcType string

const (
	PlcBool     PlcType = "bool"
	PlcByte     PlcType = "byte"
	PlcUsint    PlcType = "usint"
	PlcWord     PlcType = "word"
	PlcInt      PlcType = "int"
	PlcUint     PlcType = "uint"
	PlcDWord    PlcType = "dword"
	PlcDInt     PlcType = "dint"
	PlcUDint    PlcType = "udint"
	PlcReal     PlcType = "real"
	PlcLongReal PlcType = "longReal"
	PlcS5Time   PlcType = "s5time"
	PlcTime     PlcType = "time"
	PlcString   PlcType = "string"
)

// Configs for log file rotation and archive
type LogConf struct {
	MaxSize        int64         `yaml:"maxSize"`
	MaxAge         time.Duration `yaml:"maxAge"`
	ArchivePath    string        `yaml:"archivePath"`
	ArchiveMaxSize int64         `yaml:"archiveMaxSize"`
}

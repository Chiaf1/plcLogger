package config

import (
	"fmt"
	"os"
	"time"

	"github.com/chiaf1/plclogger/internal/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Connection struct {
		Ip   string `yaml:"ip"`
		Rack string `yaml:"rack"`
		Slot string `yaml:"slot"`
	} `yaml:"connection"`
	App struct {
		PeriodicLogInterval time.Duration `yaml:"periodicLogInterval"`
		OnChangeClock       time.Duration `yaml:"onChangeClock"`
		EnableWebServer     bool          `yaml:"enableWebServer"`
		EnableOnChangeLog   bool          `yaml:"enableOnChangeLog"`
	} `yaml:"app"`
	DataToLog []Tag `yaml:"dataToLog"`
}

type Tag struct {
	PlcTag        `yaml:",inline"`
	PeriodicLog   bool `yaml:"periodicLog"`
	OnChangeLog   bool `yaml:"onChangeLog"`
	ShowDashboard bool `yaml:"showDashboard"`
}

type PlcTag struct {
	Name    string  `yaml:"name"`
	Address string  `yaml:"address"`
	Type    PlcType `yaml:"type"`
}

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
)

// Load loads the values frrom the file "path" to the struct c, if the file is not present:
// the default values are loaded and the file is created.
func (c *Config) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.SetDefault()
			c.Save(path)
			return nil
		}
		return fmt.Errorf("Error while reading the config file: %w", err)
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("Error during parsing of YAML config file: %w", err)
	}
	return nil
}

// SetDefault sets the config default values
func (c *Config) SetDefault() {
	c.Connection.Ip = "192.168.2.2"
	c.Connection.Rack = "0"
	c.Connection.Slot = "2"

	c.App.PeriodicLogInterval = 24 * time.Hour
	c.App.PeriodicLogInterval = 30 * time.Second
	c.App.EnableWebServer = true
	c.App.EnableOnChangeLog = true

	c.DataToLog = []Tag{
		{
			PlcTag: PlcTag{
				Name:    "Tag1",
				Address: "db0.dbx0.0",
				Type:    PlcBool,
			},
			PeriodicLog:   false,
			OnChangeLog:   false,
			ShowDashboard: false,
		},
	}
}

// Save saves the configs to the "path" using the WriteFileAtomic function
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("Error while parsing to YAML: %w", err)
	}
	return utils.WriteFileAtomic(path, data, 0644)
}

// AddDataToLogFields adds the data to log to config.DataToLog slice. This metod accepts single fields.
// To use a Tag struct use the method AppendDataToLogTag()
func (c *Config) AddDataToLogFields(name string, address string, tagType PlcType, periodicLog, onChangeLog, showDashboard bool) error {
	if name == "" || address == "" || tagType == "" {
		return fmt.Errorf("missing tag values")
	}
	for _, tt := range c.DataToLog {
		if tt.Name == name {
			return fmt.Errorf("tag %s already exists", name)
		}
	}
	var newTag = Tag{
		PlcTag: PlcTag{
			Name:    name,
			Address: address,
			Type:    tagType,
		},
		PeriodicLog:   periodicLog,
		OnChangeLog:   onChangeLog,
		ShowDashboard: showDashboard,
	}
	c.DataToLog = append(c.DataToLog, newTag)
	return nil
}

// AddDataToLogTag adds the data to log to config.DataToLog slice. This metod accepts a single Tag.
// To use single fields use the method AddDataToLogFields()
func (c *Config) AppendDataToLogTag(t Tag) error {
	if t.Name == "" || t.Address == "" || t.Type == "" {
		return fmt.Errorf("missing tag values")
	}
	for _, tt := range c.DataToLog {
		if tt.Name == t.Name {
			return fmt.Errorf("tag %s already exists", t.Name)
		}
	}
	c.DataToLog = append(c.DataToLog, t)
	return nil
}

// GetPeriodicLogTags returns a slice of PlcTags that have the flag PeriodicLog enabled
func (c *Config) GetPeriodicLogTags() []PlcTag {
	var LogTags []PlcTag
	for _, t := range c.DataToLog {
		if t.PeriodicLog {
			LogTags = append(LogTags, t.PlcTag)
		}
	}
	return LogTags
}

// GetPeriodicLogTags returns a slice of PlcTags that have the flag OnChangeLog enabled
func (c *Config) GetOnChangeLogTags() []PlcTag {
	var LogTags []PlcTag
	for _, t := range c.DataToLog {
		if t.OnChangeLog {
			LogTags = append(LogTags, t.PlcTag)
		}
	}
	return LogTags
}

// GetPeriodicLogTags returns a slice of PlcTags that have the flag ShowDashboard enabled
func (c *Config) GetShowDashboardTags() []PlcTag {
	var LogTags []PlcTag
	for _, t := range c.DataToLog {
		if t.ShowDashboard {
			LogTags = append(LogTags, t.PlcTag)
		}
	}
	return LogTags
}

// RemoveDataToLog removes the tag with the "name" from the data to log slice, ritorna un errore se non la trova
func (c *Config) RemoveDataToLog(name string) error {
	for i, t := range c.DataToLog {
		if t.Name == name {
			c.DataToLog = append(c.DataToLog[:i], c.DataToLog[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Tag to remove not found")
}

// UpdateDataToLog updates the values of an existing tag
func (c *Config) UpdateDataToLogTag(t Tag) error {
	if t.Name == "" {
		return fmt.Errorf("missing tag name")
	}
	if t.Address == "" || t.Type == "" {
		return fmt.Errorf("tag %s is missing values", t.Name)
	}
	for i, tt := range c.DataToLog {
		if tt.Name == t.Name {
			c.DataToLog[i] = t
			return nil
		}
	}
	return fmt.Errorf("Tag %s to update not found", t.Name)
}

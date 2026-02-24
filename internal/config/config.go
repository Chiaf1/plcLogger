package config

import (
	"fmt"
	"os"
	"time"

	"github.com/chiaf1/plclogger/internal/domain"
	"github.com/chiaf1/plclogger/internal/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Connection domain.ConnectionConfig `yaml:"connection"`
	App        domain.AppConfig        `yaml:"app"`
	DataToLog  []domain.Tag            `yaml:"dataToLog"`
}

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
	c.Connection.Protocol = "s7"
	c.Connection.Timeout = 200 * time.Millisecond
	c.Connection.ConnectionRetry = 3

	c.App.PeriodicLogInterval = 24 * time.Hour
	c.App.OnChangeClock = 30 * time.Second
	c.App.EnableWebServer = true
	c.App.EnableOnChangeLog = true

	c.DataToLog = []domain.Tag{
		{
			PlcTag: domain.PlcTag{
				Name: "Tag1",
				Type: domain.PlcBool,
				S7: domain.S7Mapping{
					DBNumber: 0,
					Offset:   0,
					Bit:      0,
				},
				OPCUA: domain.OPCUAMapping{
					Node: "",
				},
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

// AddDataToLogTag adds the data to log to config.DataToLog slice. This metod accepts a single Tag.
// To use single fields use the method AddDataToLogFields()
func (c *Config) AppendDataToLogTag(t domain.Tag) error {
	if t.Name == "" || t.Type == "" {
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
func (c *Config) GetPeriodicLogTags() []domain.PlcTag {
	var LogTags []domain.PlcTag
	for _, t := range c.DataToLog {
		if t.PeriodicLog {
			LogTags = append(LogTags, t.PlcTag)
		}
	}
	return LogTags
}

// GetPeriodicLogTags returns a slice of PlcTags that have the flag OnChangeLog enabled
func (c *Config) GetOnChangeLogTags() []domain.PlcTag {
	var LogTags []domain.PlcTag
	for _, t := range c.DataToLog {
		if t.OnChangeLog {
			LogTags = append(LogTags, t.PlcTag)
		}
	}
	return LogTags
}

// GetPeriodicLogTags returns a slice of PlcTags that have the flag ShowDashboard enabled
func (c *Config) GetShowDashboardTags() []domain.PlcTag {
	var LogTags []domain.PlcTag
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
func (c *Config) UpdateDataToLogTag(t domain.Tag) error {
	if t.Name == "" {
		return fmt.Errorf("missing tag name")
	}
	if t.Type == "" {
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

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
		EnableWebServer     bool          `yaml:"enableWebServer"`
		EnableOnChangeLog   bool          `yaml:"enableOnChangeLog"`
	} `yaml:"app"`
	DataToLog []Tag `yaml:"dataToLog"`
}

type Tag struct {
	Name          string  `yaml:"name"`
	Address       string  `yaml:"address"`
	Type          PlcType `yaml:"type"`
	PeriodicLog   bool    `yaml:"periodicLog"`
	OnChangeLog   bool    `yaml:"onChangeLog"`
	ShowDashboard bool    `yaml:"showDashboard"`
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

func (c *Config) SetDefault() {
	c.Connection.Ip = "192.168.2.2"
	c.Connection.Rack = "0"
	c.Connection.Slot = "2"

	c.App.PeriodicLogInterval = 24 * time.Hour
	c.App.EnableWebServer = true
	c.App.EnableOnChangeLog = true

	c.DataToLog = []Tag{
		{
			Name:          "Tag1",
			Address:       "db0.dbx0.0",
			Type:          PlcBool,
			PeriodicLog:   false,
			OnChangeLog:   false,
			ShowDashboard: false,
		},
	}
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("Error while parsing to YAML: %w", err)
	}
	return utils.WriteFileAtomic(path, data, 0644)
}

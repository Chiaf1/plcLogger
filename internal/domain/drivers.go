package domain

import "github.com/robinson/gos7"

// PLC driver interface with standard metods common to multiple protocols
type PLCDriver interface {
	Connect() error
	Disconnect() error
	Read(tag PlcTag) (any, error)
	//Write(tag PlcTag) error
	//Health() bool
}

/*
#--# - S7 - #--#
*/

// struct that implements the iterface plcComunication/PLCClient
type S7Client struct {
	ip   string
	rack int
	slot int

	handler *gos7.TCPClientHandler
	client  *gos7.Client
}

type S7Mapping struct {
	DBNumber int `yaml:"dbNumber"`
	Offset   int `yaml:"offset"` // bytes offset inside the DB
	Bit      int `yaml:"bit"`    // used only for bool tags
}

/*
#--# - OPCUA - #--#
*/

// struct that implements the iterface plcComunication/PLCClient
type OPCUAClient struct {
	ip   string
	rack int
	slot int
}

type OPCUAMapping struct {
	Node string `yaml:"node"`
}

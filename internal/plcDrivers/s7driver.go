package plcdrivers

import (
	gos7 "github.com/robinson/gos7"
)

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
// NewS7Client creates the client but doesn't open the comunication
func NewS7Client(ip string, rack, slot int) *S7Client {

}
*/

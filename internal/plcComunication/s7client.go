package plccomunication

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

/*
// NewS7Client creates the client but doesn't open the comunication
func NewS7Client(ip string, rack, slot int) *S7Client {

}
*/

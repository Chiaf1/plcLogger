package plcdrivers

import "github.com/chiaf1/plclogger/internal/domain"

// struct that implements the iterface plcComunication/PLCClient
type OPCUAClient struct {
	ip   string
	rack int
	slot int
}

// NewS7Client creates the client but doesn't open the comunication
func NewOPCUAClient(conf domain.ConnectionConfig) *S7Client

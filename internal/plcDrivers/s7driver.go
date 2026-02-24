package plcdrivers

import (
	"github.com/chiaf1/plclogger/internal/domain"
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

// NewS7Client creates a new driver
func NewS7Client() *S7Client

// Connect opens the connection to the client
func (s7 *S7Client) Connect() error

// Disconnect closes the connection to the client
func (s7 *S7Client) Disconnect() error

// Read connects to the plc and reads the current value of the tag
func (s7 *S7Client) Read(tag domain.PlcTag) (any, error)

/*

	func (s7 *S7Client) Write(tag config.PlcTag) error
	func (s7 *S7Client) Health() bool
*/

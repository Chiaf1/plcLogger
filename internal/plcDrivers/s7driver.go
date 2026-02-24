package plcdrivers

import (
	"time"

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
func NewS7Client(conf domain.ConnectionConfig) *S7Client {
	h := gos7.NewTCPClientHandler(conf.Ip, conf.Rack, conf.Slot)

	if conf.Timeout <= 0 {
		conf.Timeout = 3 * time.Second
	}
	h.Timeout = conf.Timeout

	h.IdleTimeout = 5 * time.Second

	return &S7Client{
		ip:      conf.Ip,
		rack:    conf.Rack,
		slot:    conf.Slot,
		handler: h,
	}
}

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

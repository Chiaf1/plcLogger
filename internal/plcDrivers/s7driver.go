package plcdrivers

import (
	"fmt"
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
func (s7 *S7Client) Connect() error {
	if err := s7.handler.Connect(); err != nil {
		return fmt.Errorf("S7 connect failed (ip=%s rack=%s slot=%s): %w", s7.ip, s7.rack, s7.slot, err)
	}
	// after the connection let's create the client that will use this handler
	c := gos7.NewClient(s7.handler)
	s7.client = &c
	return nil
}

// Disconnect closes the connection to the client, gos7 doesn't returns errors on closing
// but to keep the interface it will return nil always
func (s7 *S7Client) Disconnect() error {
	if s7.handler != nil {
		s7.handler.Close()
	}
	return nil
}

// Read connects to the plc and reads the current value of the tag
func (s7 *S7Client) Read(tag domain.PlcTag) (any, error)

/*

	func (s7 *S7Client) Write(tag config.PlcTag) error
	func (s7 *S7Client) Health() bool
*/

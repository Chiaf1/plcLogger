package plcdrivers

import (
	"fmt"
	"sync"
	"time"

	"github.com/chiaf1/plclogger/internal/domain"
	"github.com/robinson/gos7"
)

// struct that implements the iterface plcComunication/PLCClient
type S7Client struct {
	ip   string
	rack int
	slot int

	handler *gos7.TCPClientHandler
	client  gos7.Client

	mu        sync.Mutex
	connected bool
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

// Connect opens the connection thread safe
func (s7 *S7Client) Connect() error {
	s7.mu.Lock()
	defer s7.mu.Unlock()

	if s7.connected {
		return nil // already connected
	}

	if err := s7.handler.Connect(); err != nil {
		return fmt.Errorf("S7 connect failed (ip=%s rack=%v slot=%v): %w", s7.ip, s7.rack, s7.slot, err)
	}
	// after the connection let's create the client that will use this handler
	s7.client = gos7.NewClient(s7.handler)

	s7.connected = true
	return nil
}

// Disconnect closes the connection to the client, gos7 doesn't returns errors on closing
// but to keep the interface it will return nil always
func (s7 *S7Client) Disconnect() error {
	s7.mu.Lock()
	defer s7.mu.Unlock()

	if s7.connected && s7.handler != nil {
		s7.handler.Close()
	}
	s7.client = nil
	s7.connected = false
	return nil
}

// Read connects to the plc and reads the current value of the tag, if it loses connection it retrys 2 times with a time out in between
func (s7 *S7Client) Read(tag domain.PlcTag) (any, error) {
	const maxAttempts = 2

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// check connection
		if err := s7.Connect(); err != nil {
			if attempt == maxAttempts {
				return nil, fmt.Errorf("PLC connection failed: %w", err)
			}
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// read attempt
		val, err := s7.readOnce(tag)
		if err == nil {
			return val, nil
		}

		// if it fails we try to reconnect
		s7.Disconnect()

		if attempt == maxAttempts {
			return nil, fmt.Errorf("read failed after max attempts: %w", err)
		}

		time.Sleep(200 * time.Millisecond)
	}

	return nil, fmt.Errorf("unexpected read failure")
}

// readOnce reads from a db the single tag in the argument
func (s7 *S7Client) readOnce(tag domain.PlcTag) (any, error) {
	// I make a copy of the pointer to the s7Client so I can work on it in parallel with other functions
	s7.mu.Lock()
	client := s7.client
	s7.mu.Unlock()

	// Calculate how many bytes are needed
	sz, err := requiredSizeBytes(tag)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, sz)

	// Reading from DB
	err = client.AGReadDB(tag.S7.DBNumber, tag.S7.Offset, sz, buf)
	if err != nil {
		return nil, fmt.Errorf("error reading from (DB=%v offset=%v, size=%v): %w", tag.S7.DBNumber, tag.S7.Offset, sz, err)
	}
	v, err := parseValueFromBuffer(buf, 0, tag)
	if err != nil {
		return nil, fmt.Errorf("error parsing the byte buffer: %w", err)
	}
	return v, nil
}

// Calculates how many bytes are needed for a given tag
func requiredSizeBytes(tag domain.PlcTag) (int, error) {
	switch tag.Type {
	case domain.PlcBool:
		return 1, nil
	case domain.PlcByte, domain.PlcUsint:
		return 1, nil
	case domain.PlcWord:
		return 2, nil
	case domain.PlcInt:
		return 2, nil
	case domain.PlcUint:
		return 2, nil
	case domain.PlcDWord:
		return 4, nil
	case domain.PlcDInt:
		return 4, nil
	case domain.PlcUDint:
		return 4, nil
	case domain.PlcReal:
		return 4, nil
	case domain.PlcLongReal:
		return 8, nil
	case domain.PlcS5Time:
		return 2, nil
	case domain.PlcTime:
		return 4, nil
	case domain.PlcString:
		return 255, nil
	default:
		return -1, fmt.Errorf("unsupported type: %s", tag.Type)
	}
}

// parseValueFromBuffer parses the tag value from the read buffer
func parseValueFromBuffer(buf []byte, startIndex int, tag domain.PlcTag) (any, error) {
	if startIndex < 0 || startIndex >= len(buf) {
		return nil, fmt.Errorf("StartIndex out of range: %d (len=%d)", startIndex, len(buf))
	}
	var h gos7.Helper
	switch tag.Type {
	case domain.PlcBool:
		return h.GetBoolAt(buf[startIndex], uint(tag.S7.Bit)), nil
	case domain.PlcByte, domain.PlcUsint:
		return buf[startIndex], nil
	case domain.PlcWord:
		return getWordAt(buf, startIndex), nil
	case domain.PlcInt:
		return getIntAt(buf, startIndex), nil
	case domain.PlcUint:
		return getUintAt(buf, startIndex), nil
	case domain.PlcDWord:
		return getDwordAt(buf, startIndex), nil
	case domain.PlcDInt:
		return getDintAt(buf, startIndex), nil
	case domain.PlcUDint:
		return getUDintAt(buf, startIndex), nil
	case domain.PlcReal:
		return h.GetRealAt(buf, startIndex), nil
	case domain.PlcLongReal:
		return h.GetLRealAt(buf, startIndex), nil
	case domain.PlcS5Time:
		return h.GetS5TimeAt(buf, startIndex), nil
	case domain.PlcTime:
		return getDintAt(buf, startIndex), nil
	case domain.PlcString:
		return h.GetStringAt(buf, startIndex), nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", tag.Type)
	}
}

// helpers functions for converting byte buffers into values

// returns the uint at index pos swapping the bytes
func getUintAt(buf []byte, pos int) uint16 {
	return uint16(buf[pos])<<8 | uint16(buf[pos+1])
}

// returns int at index pos swaping the bytes
func getIntAt(buf []byte, pos int) int16 {
	return int16(getUintAt(buf, pos))
}

func getWordAt(buf []byte, pos int) uint16 {
	return getUintAt(buf, pos)
}

func getDwordAt(buf []byte, pos int) uint32 {
	return uint32(buf[pos])<<24 | uint32(buf[pos+1])<<16 | uint32(buf[pos+2])<<8 | uint32(buf[pos+3])
}

func getUDintAt(buf []byte, pos int) uint32 {
	return getDwordAt(buf, pos)
}

func getDintAt(buf []byte, pos int) int32 {
	return int32(getDwordAt(buf, pos))
}

/*
	func (s7 *S7Client) Write(tag config.PlcTag) error
	func (s7 *S7Client) Health() bool
*/

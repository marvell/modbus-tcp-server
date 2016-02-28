package modbus

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type Request struct {
	ADU     *ADU
	addr    string
	t       time.Time
	Context map[string]interface{}
}

func readRequest(c net.Conn) (*Request, error) {
	buf := make([]byte, 512)
	n, err := c.Read(buf)
	if err != nil {
		return nil, err
	}
	buf = buf[:n]

	if len(buf) < 8 {
		return nil, fmt.Errorf("wrong format")
	}

	adu := &ADU{
		ID:       binary.BigEndian.Uint16(buf[0:2]),
		Proto:    binary.BigEndian.Uint16(buf[2:4]),
		Length:   binary.BigEndian.Uint16(buf[4:6]),
		Device:   uint8(buf[6]),
		Function: uint8(buf[7]),
		Data:     buf[8:],
	}

	return &Request{
		adu,
		c.RemoteAddr().String(),
		time.Now(),
		make(map[string]interface{}),
	}, nil
}

func (r *Request) RemoteAddr() string {
	return r.addr
}

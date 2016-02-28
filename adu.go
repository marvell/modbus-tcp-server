package modbus

import (
	"encoding/binary"
)

type ADU struct {
	ID       uint16
	Proto    uint16
	Length   uint16
	Device   uint8
	Function uint8
	Data     []byte
}

func (a ADU) Bytes() []byte {
	b := make([]byte, 8)

	binary.BigEndian.PutUint16(b[0:2], a.ID)
	binary.BigEndian.PutUint16(b[2:4], a.Proto)
	binary.BigEndian.PutUint16(b[4:6], uint16(2+len(a.Data)))
	b[6] = a.Device
	b[7] = a.Function
	b = append(b, a.Data...)

	return b
}

package modbus

import "fmt"

type Error uint8

var (
	ErrIllegalFunc     Error = 1
	ErrIllegalDataAddr Error = 2
	ErrIllegalDataVal  Error = 3
	ErrDeviceFailure   Error = 4
)

func (e Error) Error() string {
	return fmt.Sprintf("%d", e)
}

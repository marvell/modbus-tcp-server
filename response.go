package modbus

type Response struct {
	req  *Request
	data []byte
	err  *Error
}

func (r *Response) Bytes() []byte {
	adu := ADU{
		ID:       r.req.ADU.ID,
		Proto:    r.req.ADU.Proto,
		Length:   6,
		Device:   r.req.ADU.Device,
		Function: r.req.ADU.Function,
		Data:     r.data,
	}

	if r.err != nil {
		adu.Function = adu.Function | 128
		adu.Data = []byte{byte(*r.err)}
	}

	return adu.Bytes()
}

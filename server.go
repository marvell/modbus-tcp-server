package modbus

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

var (
	Debug bool = false
)

func NewServer() *Server {
	return &Server{}
}

type Server struct {
	pre func(req *Request)
	f16 func(req *Request, start, count int, regs [][2]byte) *Error
}

func (s *Server) RegisterPreHandler(f func(req *Request)) {
	s.pre = f
}

func (s *Server) RegisterHandler(n int, f interface{}) error {
	switch n {
	case 16:
		s.f16 = f.(func(req *Request, start, count int, regs [][2]byte) *Error)
		return nil
	}

	return errors.New("wrong number of function")
}

func (s *Server) ListenOn(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go func(c net.Conn) {
			defer c.Close()

			req, err := readRequest(c)
			if Debug {
				log.Printf("req: %#v", req)
			}
			if err != nil {
				if err != io.EOF {
					log.Printf("error while reading request from %s: %s", c.RemoteAddr(), err)
				}

				return
			}

			if s.pre != nil {
				s.pre(req)
			}

			res, err := s.handle(req)
			if Debug {
				log.Printf("res: %#v", res)
			}
			if err != nil {
				log.Printf("response error: %s", err)
				return
			}

			c.Write(res.Bytes())
		}(conn)
	}

	return nil
}

func (server *Server) handle(req *Request) (*Response, error) {
	var err *Error
	var data []byte

	switch {
	case req.ADU.Function == 0x10 && server.f16 != nil:
		s := int(binary.BigEndian.Uint16(req.ADU.Data[0:2]))
		c := int(binary.BigEndian.Uint16(req.ADU.Data[2:4]))
		d := req.ADU.Data[5:]

		if len(d)/2 != c {
			err = &ErrIllegalDataAddr
			break
		}

		regs := make([][2]byte, c)
		for i := 0; i < int(c); i++ {
			regs[i] = [2]byte{d[2*i], d[2*i+1]}
		}

		err = server.f16(req, s, c, regs)
		if err == nil {
			data = req.ADU.Data[0:4]
		}
	default:
		err = &ErrIllegalFunc
	}

	return &Response{req, data, err}, nil
}

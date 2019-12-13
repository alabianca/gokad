package gokad

import (
	"encoding/binary"
	"net"
)

type Contact struct {
	ID   ID
	IP   net.IP
	Port int
	next *Contact
}

func (c *Contact) Serialize() ([]byte, error) {
	id := c.ID
	ip, err := c.IP.MarshalText()
	if err != nil {
		return nil, err
	}

	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(c.Port))

	end := []byte("/")

	concat := make([]byte, 0)
	concat = append(concat, id...)
	concat = append(concat, port...)
	concat = append(concat, ip...)
	concat = append(concat, end...)

	return concat, nil

}

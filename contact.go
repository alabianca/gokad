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

// 20 bytes id <- 16 bytes ip <- 2 bytes port <- 1 byte end
func (c *Contact) Serialize() []byte {
	id := c.ID

	ip := []byte(c.IP)

	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(c.Port))

	end := make([]byte, 0)

	concat := make([]byte, 0)
	concat = append(concat, id...)
	concat = append(concat, ip...)
	concat = append(concat, port...)
	concat = append(concat, end...)

	return concat

}

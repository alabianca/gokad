package gokad

import "net"

type Contact struct {
	ID   *ID
	IP   net.IP
	Port int
	Next *Contact
}

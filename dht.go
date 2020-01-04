package gokad

import (
	"net"
)

const MessageSize = 800

const k = 20

type DHT struct {
	ID           ID
	RoutingTable *RoutingTable
}

func NewDHT() *DHT {
	id := GenerateRandomID()
	routing := NewRoutingTable(id)

	return &DHT{
		ID:           id,
		RoutingTable: routing,
	}
}

func (dht *DHT) Bootstrap(port int, ip net.IP, hexId string) (Contact, int, error) {
	id, err := From(hexId)
	if err != nil {
		return Contact{}, 0, err
	}
	c := Contact{
		IP:   ip,
		Port: port,
		ID:   id,
	}

	return dht.RoutingTable.Add(c)
}

func (dht *DHT) GetAlphaNodes(alpha int, id ID) []Contact {
	return dht.RoutingTable.GetAlphaNodes(alpha, id)
}

// RPC
func (dht *DHT) FindNode(id ID) []Contact {
	return dht.GetAlphaNodes(k, id)
}

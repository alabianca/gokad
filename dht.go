package gokad

import "net"

const MessageSize = 800

const k = 20

type Value struct{
	Host net.IP
	Port int
}

type values map[string]Value

type DHT struct {
	ID           ID
	routingTable *RoutingTable
	storedValues values
}

func NewDHT() *DHT {
	id := GenerateRandomID()
	routing := NewRoutingTable(id)

	return &DHT{
		ID:           id,
		routingTable: routing,
		storedValues: make(values),

	}
}

func (dht *DHT) RoutingTable() *RoutingTable {
	return dht.routingTable
}

func (dht *DHT) GetAlphaNodes(alpha int, id ID) []Contact {
	return dht.routingTable.GetAlphaNodes(alpha, id)
}

// RPC
func (dht *DHT) FindNode(id ID) []Contact {
	return dht.GetAlphaNodes(k, id)
}

func (dht *DHT) Store(key ID, ip net.IP, port int) {
	dht.storedValues[key.String()] = Value{
		Host: ip,
		Port: port,
	}
}

func (dht *DHT) FindValue(key ID) ([]Contact, Value) {
	v, ok := dht.storedValues[key.String()]
	if ok {
		return nil, v
	}

	return dht.GetAlphaNodes(k, key), Value{}
}

package gokad

type DHT struct {
	ID           *ID
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

func (dht *DHT) Bootstrap(bootsrapNode *Contact) {

}

// func (dht *DHT) RPCNodeLookup(hex string) error {
// 	target, err := From(hex)
// 	if err != nil {
// 		return err
// 	}

// 	lookupNodes := dht.RoutingTable.getXClosestContacts(3, target)

// 	return nil
// }

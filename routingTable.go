package gokad

import "math"

// RoutingTable that hold the KBuckets
type RoutingTable struct {
	id      *ID
	buckets []*KBucket
}

// NewRoutingTable returns a newly ininitalized routing table
// The routing table's size is determined by MaxRoutingTableSize which is set to 160
func NewRoutingTable(id *ID) *RoutingTable {
	r := &RoutingTable{
		id:      id,
		buckets: make([]*KBucket, MaxRoutingTableSize),
	}

	for i := range r.buckets {
		r.buckets[i] = NewKBucket(i)
	}

	return r

}

// Add adds a new contact into the appropriate k-bucket within the routing table
// returning the head of the bucket, the insertion index and an error if there was one
func (r *RoutingTable) Add(c *Contact) (*Contact, int, error) {
	delta := r.id.DistanceTo(c.ID)
	index := r.determineInsertIndex(delta)

	head, err := r.insertAt(index, c)

	return head, index, err
}

// GetAlphaNodes gets α nodes out of its k-bucket where the id to be looked up would fit in.
// α is a system wide concurrency parameter a value of 3 is suggested. If the corresponding k-bucket
// has less than α entries, the node takes the α closest nodes it knows of.
// Source: Implementation of the Kademlia Hash Table by Bruno Spori
// https://pub.tik.ee.ethz.ch/students/2006-So/SA-2006-19.pdf
func (r *RoutingTable) GetAlphaNodes(alpha int, id *ID) []Contact {
	return r.getXClosestContacts(alpha, id)
}

func (r *RoutingTable) insertAt(i int, c *Contact) (*Contact, error) {
	bucket := r.buckets[i]
	return bucket.Insert(c)
}

// determineInsertIndex determines at which index the id/contact should be inserted based on the distance
func (r *RoutingTable) determineInsertIndex(delta NodeID) int {

	for i := 0; i < MaxRoutingTableSize; i++ {
		bit := delta.GetBitAt(uint(i))
		if bit > 0 {
			return MaxRoutingTableSize - 1 - i
		}
	}

	return 0

}

// determineOrderOfVisits determines in which order we need to visit
// k-buckets to close node to delta.
//
// For example assume that the own id is 1001110110000101 (Node1) and the id to which the closest
// ids shall be searched is 1001111110100001 (Node2). To get the distance between the two, they have to be
// xored, to get a delta of 0000001000100100. Using this delta we know that the closest nodes to Node2
// must be in k-bucket number 9. We then visit k-bucket number 5 and then number 2 and then number 0 and
// then go ascending visiting each zero bit bucket
func (r *RoutingTable) determineOrderOfVisits(delta NodeID) []int {
	phase1 := make([]int, 0)
	phase2 := make([]int, 0)

	// loop through the entire address space.
	// ie 0 -> 159
	for i := 0; i < MaxRoutingTableSize; i++ {
		bit := delta.GetBitAt(uint(i))
		realBucketIndex := bucketIndex(i)
		if bit > 0 {
			phase1 = append(phase1, realBucketIndex)
		} else {
			phase2 = append([]int{realBucketIndex}, phase2...)

		}
	}

	phase1 = append(phase1, phase2...)
	return phase1
}

func (r *RoutingTable) getXClosestContacts(x int, id *ID) []Contact {
	delta := r.id.DistanceTo(id)
	visits := r.determineOrderOfVisits(delta)
	out := make([]Contact, 0)

	for _, bucketIndex := range visits {
		bucket := r.buckets[bucketIndex]
		res := bucket.getXClosestContacts(x, id)

		if len(res) > 0 {
			for _, c := range res {
				out = append(out, *c)

			}
		}

		if len(out) >= x {
			break
		}
	}

	bounds := int(math.Min(float64(x), float64(len(out))))

	return out[:bounds]

}

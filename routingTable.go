package gokad

import "math"

// routingTable that hold the KBuckets
type RoutingTable struct {
	id      ID
	buckets []*KBucket
}

// NewRoutingTable returns a newly ininitalized routing table
// The routing table's size is determined by MaxRoutingTableSize which is set to 160
func NewRoutingTable(id ID) *RoutingTable {
	r := &RoutingTable{
		id:      id,
		buckets: make([]*KBucket, MaxRoutingTableSize),
	}

	for i := range r.buckets {
		r.buckets[i] = NewKBucket(i)
	}

	return r

}

/*  Add adds a new contact into the appropriate k-bucket within the routing table
    returning the contact that was added OR the head of the bucket, the insertion index and an error if there was one.
    The head of the bucket is only returned if there is also a Bucket ErrBucketAtCapacity error.
    We do this so we can ping the head to see if it is still active

        "If Bucket contains MaxCapacity, the node at the head is pinged. If it replies, the current head is moved
        to the tail and the contact is not added. If it does not reply, the head is discarded and the contact is
        added to the tail"
	@source: Implementation of the Kademlia Distributed Hash Table by Bruno Spori Semester Thesis
   https://pub.tik.ee.ethz.ch/students/2006-So/SA-2006-19.pdf
**/
func (r *RoutingTable) Add(c Contact) (Contact, int, error) {
	delta := r.id.DistanceTo(c.ID)
	index := r.determineBucketIndex(delta)

	contactOrHead, err := r.insertAt(index, c)

	return contactOrHead, index, err
}

func (r *RoutingTable) Bucket(index int) (KBucket, bool) {
	if len(r.buckets) == 0 {
		return KBucket{}, false
	}
	if index > len(r.buckets) {
		return *r.buckets[len(r.buckets) - 1], true
	}

	if index < 0 {
		return *r.buckets[0], true
	}

	return *r.buckets[index], true
}

// GetAlphaNodes gets α nodes out of its k-bucket where the id to be looked up would fit in.
// α is a system wide concurrency parameter a value of 3 is suggested. If the corresponding k-bucket
// has less than α entries, the node takes the α closest nodes it knows of.
// Source: Implementation of the Kademlia Hash Table by Bruno Spori
// https://pub.tik.ee.ethz.ch/students/2006-So/SA-2006-19.pdf
func (r *RoutingTable) GetAlphaNodes(alpha int, id ID) []Contact {
	return r.getXClosestContacts(alpha, id)
}

func (r *RoutingTable) insertAt(i int, c Contact) (Contact, error) {
	bucket := r.buckets[i]
	return bucket.Insert(c)
}

// determineBucketIndex determines at which index the id/contact should be inserted based on the distance
func (r *RoutingTable) determineBucketIndex(delta NodeID) int {

	for i := 0; i < MaxRoutingTableSize; i++ {
		bit := delta.GetBitAt(uint(i))
		if bit > 0 {
			return MaxRoutingTableSize - 1 - i
		}
	}

	return 0

}

// determineOrderOfVisits determines in which order we need to visit
// our k-buckets to find nodes close to delta
// we do so by looping through the bit string of delta.
// consider a delta of 1 0 1 0 (in an address space of 4 for brevity)
// The order in which we visit the k-buckets will be:
//   bucket index 3 -> bucket index 1 -> bucket index 0 -> bucket index 2
// It may be better explained by looking the the full address space as a binary tree
/**
                                 /         \
                                /           \
                               /             \
                              /               \
                             /                 \
                            /                   \
                           /                     \
                          / \                   / \
                         /   \                 /   \
                        /     \               /     \
                       /       \             /       \
                      /         \           /         \
                     /           \         /           \
                    / \         / \       /\           /\
                   /   \       /   \     /  \         /  \
                  /\   /\     /\  / \   / \ /\       /\ / \
                 0 1  2 3    4 5 6  7  8  91011     12131415

    Consider Our ID: 0010 (2)
             Target: 1000 (8)
                    ---------
             Delta:  1010

   The first bucket we visit is bucket 3. This bucket contains ids ranging 8 - 15.
   In this bucket the deltas range from 0000 - 0111
   The next bucket we visit is bucket 1. This bucket contains ids ranging from 0 - 1
   In this bucket the deltas range from 1000 - 1001
   The next bucket we visit is bucket 0. This bucket contains ids ranging from 2 - 3
   In this bucket the deltas range from 1010 - 1011
   The next and last bucket we visit is bucket 2. This bucket contains ids ranging from 4 - 7
   In this bucket the deltas range from 1100 - 1111

As you can see with every bucket we visit the delta to our Target increases.

 */
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

func (r *RoutingTable) getXClosestContacts(x int, id ID) []Contact {
	delta := r.id.DistanceTo(id)
	visits := r.determineOrderOfVisits(delta)
	out := make([]Contact, 0)

	for _, bucketIndex := range visits {
		bucket := r.buckets[bucketIndex]
		res := bucket.getXClosestContacts(x, id)

		if len(res) > 0 {
			for _, c := range res {
				out = append(out, c)

			}
		}

		if len(out) >= x {
			break
		}
	}

	bounds := int(math.Min(float64(x), float64(len(out))))

	return out[:bounds]

}

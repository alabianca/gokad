package gokad

import (
	"bytes"
	"errors"
	"math"
)

// MaxCapacity is a system defined MaxCapacity of each kbucket
const MaxCapacity = 20

// Errors
const ErrBucketAtCapacity = "Bucket at Capacity"
const ErrContactExists = "Contact Exists Already"
const ErrBucketIndexOutOfBounds = "Bucket Index Out Of Bounds"
const ErrNoHeadFound = "No Bucket Head Found"

// KBucket is a bucket that contains k (MaxCapacity) contacts
type KBucket struct {
	Index int
	head  *Contact
	tail  *Contact
	size  int
}

// NewKBucket returns a new KBucket with Index index
func NewKBucket(index int) *KBucket {
	return &KBucket{
		Index: index,
	}
}

// Insert attempts to insert a new Contact into the KBucket
// Adding a new node to the bucket contains the following steps:
//  1. If Bucket contains less than MaxCapacity nodes and node does not already exist - add node to tail
//  2. If Bucket contains node already, the node is moved to the tail of the list
//  3. If Bucket contains MaxCapacity, the node at the head is pinged. If it replies, the current head is moved
//     to the tail and the contact is not added. If it does not reply, the head is discarded and the contact is
//     added to the tail
// @Source: Implementation of the Kademlia Distributed Hash Table by Bruno Spori Semester Thesis
// https://pub.tik.ee.ethz.ch/students/2006-So/SA-2006-19.pdf
///
func (b *KBucket) Insert(c *Contact) (*Contact, error) {
	// bucket is completely empty. just initialize it
	if b.IsEmpty() {
		b.add(c)
		return c, nil
	}

	index := b.indexOf(c)
	// 2. Node already exists: Move the node to the tail
	if index > -1 {
		b.moveToTail(index)
		return c, errors.New(ErrContactExists)
		// 1. Bucket does not contain node and is not at capacity: add it to the tail
	} else if index < 0 && b.size < MaxCapacity {
		b.add(c)
		return c, nil
	}

	return b.head, errors.New(ErrBucketAtCapacity)

}

// Walk traverses the bucket list calling the walkFn for each contact in the bucket
// if need to return early from the walk, return true from the walkFn
func (b *KBucket) Walk(walkFn func(c *Contact) bool) {
	head := b.head
	if head == nil {
		return
	}

	for {
		done := walkFn(head)
		head = head.next
		if head == nil || done {
			break
		}
	}
}

// IsEmpty returns true if the bucket is empty. false otherwise
func (b *KBucket) IsEmpty() bool {
	return b.head == nil && b.tail == nil
}

// Size returns the size of the bucket
func (b *KBucket) Size() int {
	return b.size
}

// Head returns the head of the bucket list's contacts
func (b *KBucket) Head() *Contact {
	return b.head
}

// Tail returns the tail of the bucket list's contacts
func (b *KBucket) Tail() *Contact {
	return b.tail
}

// indexOf returns the index of the contact
// if the contact is not found it returns -1
func (b *KBucket) indexOf(c *Contact) int {
	index := -1
	var found bool
	b.Walk(func(contact *Contact) bool {
		index++
		if c.ID.Equal(contact.ID) {
			found = true
			return true
		}

		return false
	})

	if !found {
		return -1
	}

	return index
}

func (b *KBucket) add(c *Contact) {
	b.size++
	if b.IsEmpty() {
		b.head = c
		b.tail = c
	} else {
		b.tail.next = c
		b.tail = c
	}

}

func (b *KBucket) moveToTail(index int) error {
	if index >= b.size {
		return errors.New(ErrBucketIndexOutOfBounds)
	}

	head := b.head
	if head == nil {
		return errors.New(ErrNoHeadFound)
	}

	if head.next == nil || head.next.next == nil {
		return nil
	}

	if index == 0 {
		temp := b.head
		b.head = b.head.next
		b.tail = temp
		return nil
	}

	counter := 1
	var slow *Contact
	var target *Contact
	var fast *Contact

	slow = head
	target = head.next
	fast = target.next

	for {
		if counter == index {
			slow.next = fast
			target.next = nil
			b.tail = target
			return nil
		}

		counter++
		slow = target
		target = fast
		fast = fast.next

		if fast == nil && counter != index {
			return errors.New(ErrBucketIndexOutOfBounds)
		}
	}

}

func (b *KBucket) getXClosestContacts(x int, targetID ID) []*Contact {
	distances := make([]Distance, b.Size())
	index := 0
	distanceMap := make(map[string]*Contact)

	b.Walk(func(c *Contact) bool {
		delta := c.ID.DistanceTo(targetID)
		distanceMap[delta.String()] = c
		distances[index] = delta
		index++
		return false
	})

	// now sort the list by distance
	sort(distances)
	bounds := math.Min(float64(x), float64(len(distances)))
	xClosestDeltas := distances[:int(bounds)]
	out := make([]*Contact, 0)

	for _, d := range xClosestDeltas {
		contact, ok := distanceMap[d.String()]
		if ok {
			out = append(out, contact)
		}
	}

	return out
}

func (b KBucket) String() string {
	buf := new(bytes.Buffer)

	b.Walk(func(c *Contact) bool {
		buf.WriteString(c.ID.String() + ",")
		return false
	})

	return buf.String()
}

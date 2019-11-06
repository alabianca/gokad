package gokad

import (
	"log"
	"math/rand"
	"net"
	"testing"
	"time"
)

func init() {
	log.Println("Seeding Random")
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestAddToBucket(t *testing.T) {
	bucket := NewKBucket(0)
	c1 := generateRandomContact()
	c2 := generateRandomContact()
	c3 := generateRandomContact()
	c4 := generateRandomContact()
	bucket.add(c1)
	bucket.add(c2)
	bucket.add(c3)
	bucket.add(c4)

	if bucket.Size() != 4 {
		t.Errorf("Expected size to be %d but got %d\n", 4, bucket.Size())
	}

	if !bucket.Head().ID.Equal(c1.ID) {
		t.Errorf("Expected head to be %s but got %s", bucket.Head().ID.String(), c1.ID.String())
	}

	if !bucket.Tail().ID.Equal(c4.ID) {
		t.Errorf("Expected tail to be %s but got %s", bucket.Tail().ID.String(), c4.ID.String())
	}

	str := bucket.String()
	expect := c1.ID.String() + "," + c2.ID.String() + "," + c3.ID.String() + "," + c4.ID.String() + ","

	if str != expect {
		t.Errorf("Expected %s\nbut got\n%s", str, expect)
	}

}

func TestMoveHeadToTail(t *testing.T) {
	bucket := NewKBucket(0)
	c1 := generateRandomContact()
	c2 := generateRandomContact()
	c3 := generateRandomContact()
	c4 := generateRandomContact()
	bucket.add(c1)
	bucket.add(c2)
	bucket.add(c3)
	bucket.add(c4)

	t.Logf("c1 %s\n", c1.ID.String())
	t.Logf("c2 %s\n", c2.ID.String())
	t.Logf("c3 %s\n", c3.ID.String())
	t.Logf("c4 %s\n", c4.ID.String())

	bucket.moveToTail(0)

	if !bucket.Tail().ID.Equal(c1.ID) {
		t.Errorf("Expected Tail to be %s but got %s\n", c1.ID.String(), bucket.Tail().ID.String())
	}

	if !bucket.Head().ID.Equal(c2.ID) {
		t.Errorf("Expected Head to be %s but got %s\n", c2.ID.String(), bucket.Head().ID.String())
	}

}

func TestMoveNToTail(t *testing.T) {
	bucket := NewKBucket(0)
	c1 := generateRandomContact()
	c2 := generateRandomContact()
	c3 := generateRandomContact()
	c4 := generateRandomContact()
	bucket.add(c1)
	bucket.add(c2)
	bucket.add(c3)
	bucket.add(c4)

	t.Logf("c1 %s\n", c1.ID.String())
	t.Logf("c2 %s\n", c2.ID.String())
	t.Logf("c3 %s\n", c3.ID.String())
	t.Logf("c4 %s\n", c4.ID.String())

	bucket.moveToTail(2)

	if !bucket.Tail().ID.Equal(c3.ID) {
		t.Errorf("Expected Tail to be %s but got %s\n", c3.ID.String(), bucket.Tail().ID.String())
	}

	if !bucket.Head().ID.Equal(c1.ID) {
		t.Errorf("Expected Head to be %s but got %s\n", c1.ID.String(), bucket.Head().ID.String())
	}

}

func TestGetXclosestContacts(t *testing.T) {
	bucket := NewKBucket(0)
	//c1 := generateContactFrom("C80F741BC1B397C54A54858E4E2A8840B2BC766B")
	c1 := generateContactFrom("D80F741BC1B397C54A54858E4E2A8840B2BC766B")
	c2 := generateContactFrom("F70F741BC1B397C54A54858E4E2A8840B2BC766B")
	c3 := generateContactFrom("A70F441BC1B397C54A54858E4E2A8840B2BC766B")

	bucket.add(c3)
	bucket.add(c1)
	bucket.add(c2)

	targetID, _ := From("C80F741BC1B397C54A54858E4E2A8840B2BC766B")

	out := bucket.getXClosestContacts(2, targetID)

	if len(out) != 2 {
		t.Errorf("Expected length to be %d but got %d\n", 2, len(out))
	}

	for i, c := range out {
		if i == 0 {
			if c.ID.String() != c1.ID.String() {
				t.Errorf("Expected (i0) to be %s but got %s\n", c1.ID, c.ID)

			}

			continue
		}

		if i == 1 {
			if c.ID.String() != c2.ID.String() {
				t.Errorf("Expected (i1) to be %s but got %s\n", c2.ID, c.ID)
			}

			continue
		}
	}

}

// Utils

func getPreSetBucket() *KBucket {
	bucket := NewKBucket(0)
	c1 := generateRandomContact()
	c2 := generateRandomContact()
	c3 := generateRandomContact()
	c4 := generateRandomContact()
	bucket.add(c1)
	bucket.add(c2)
	bucket.add(c3)
	bucket.add(c4)

	return bucket
}

func generateRandomContact() *Contact {
	id := GenerateRandomID()
	return &Contact{
		ID:   id,
		IP:   net.IPv4(byte(127), byte(0), byte(0), byte(1)), // 127.0.0.1
		Port: 3000,
	}
}

func generateContactFrom(hex string) *Contact {
	id, _ := From(hex)
	return &Contact{
		ID:   id,
		IP:   net.IPv4(byte(127), byte(0), byte(0), byte(1)), // 127.0.0.1
		Port: 3000,
	}
}

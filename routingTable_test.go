package gokad

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestInsertDistanceBetween2Ids(t *testing.T) {

	cases := []struct {
		id1 string
		id2 string
		out int
	}{
		{
			id1: "480F741BC1B397C54A54858E4E2A8840B2BC766B", // 0100100000001111011101000001101111000001101100111001011111000101010010100101010010000101100011100100111000101010100010000100000010110010101111000111011001101011
			id2: "C80F741BC1B397C54A54858E4E2A8840B2BC766B", // 1100100000001111011101000001101111000001101100111001011111000101010010100101010010000101100011100100111000101010100010000100000010110010101111000111011001101011
			out: 159,
		},
		{
			id1: "480F741BC1B397C54A54858E4E2A8840B2BC766B",
			id2: "480F741BC1B397C54A54858E4E2A8840B2BC766B",
			out: 0,
		},
		{
			id1: "489887A2C81C7920911815BCD99D3F19AB3D633D",
			id2: "DC84DABE02FB31B9011800635C03794213EBA1F0",
			out: 159,
		},
	}

	for _, c := range cases {
		id1, _ := From(c.id1)
		id2, _ := From(c.id2)
		routing := NewRoutingTable(id1)
		delta := id1.DistanceTo(id2)
		index := routing.determineBucketIndex(delta)
		if index != c.out {
			t.Logf("Id1: %s\nId2: %s\n", c.id1, c.id2)
			t.Logf("Delta: %x\n", delta)
			t.Errorf("Expected index to be %d, but got %d", c.out, index)
		}
	}

}

func TestDetermineInsertIndex(t *testing.T) {
	cases := []struct {
		IN  string
		OUT int
	}{
		{
			IN:  "24d48bcb03a8b24aaf531c99bef39bee2d666496",
			OUT: 157,
		},
		{
			IN:  "858b07598756998c87f48cdafdd9fad9f1a714ff",
			OUT: 159,
		},
	}

	for _, c := range cases {
		id, _ := From(c.IN)
		routing := NewRoutingTable(nil)
		index := routing.determineBucketIndex(id)
		if index != c.OUT {
			t.Errorf("Expected %d but got %d\n", c.OUT, index)
		}
	}

}

func TestAddContactToRoutingTableWithoutErrors(t *testing.T) {
	cases := []struct {
		id1 string
		id2 string
		out int
	}{
		{
			id1: "480F741BC1B397C54A54858E4E2A8840B2BC766B", // 0100100000001111011101000001101111000001101100111001011111000101010010100101010010000101100011100100111000101010100010000100000010110010101111000111011001101011
			id2: "C80F741BC1B397C54A54858E4E2A8840B2BC766B", // 1100100000001111011101000001101111000001101100111001011111000101010010100101010010000101100011100100111000101010100010000100000010110010101111000111011001101011
			out: 159,
		},
		{
			id1: "480F741BC1B397C54A54858E4E2A8840B2BC766B",
			id2: "480F741BC1B397C54A54858E4E2A8840B2BC766B",
			out: 0,
		},
		{
			id1: "489887A2C81C7920911815BCD99D3F19AB3D633D",
			id2: "DC84DABE02FB31B9011800635C03794213EBA1F0",
			out: 159,
		},
	}

	for _, c := range cases {
		id, _ := From(c.id1)
		contact := generateContactFrom(c.id2)
		routing := NewRoutingTable(id)

		addedC, i, err := routing.Add(contact)

		if err != nil {
			t.Errorf("Expected error to be nil but got: %s\n", err)
		}

		if i != c.out {
			t.Errorf("Expected insert index to be %d but got %d\n", c.out, i)
		}

		if strings.ToLower(addedC.ID.String()) != strings.ToLower(c.id2) {
			t.Errorf("Expected added contact to be %s but got %s\n", c.id2, addedC.ID)
		}

		if routing.buckets[i].Size() != 1 {
			t.Errorf("Expected bucket size to be 1, but got %d\n", routing.buckets[i].Size())
		}
	}
}

func TestAddExistingContactToRoutingTable(t *testing.T) {
	contact1 := generateContactFrom("480F741BC1B397C54A54858E4E2A8840B2BC766B")
	contact2 := generateContactFrom("480F741BC1B397C54A54858E4E2A8840B2BC766B")
	id, _ := From("489887A2C81C7920911815BCD99D3F19AB3D633D")
	routing := NewRoutingTable(id)

	// Add first contact. Expect it to be added without issues
	head, _, err := routing.Add(contact1)
	if err != nil {
		t.Errorf("Expected error to be nil but got: %s\n", err)
	}

	// Add contact again. Should be not added again
	head, _, err = routing.Add(contact2)
	if reflect.DeepEqual(head, Contact{}) {
		t.Errorf("Expected ping head to be %s but got empty contact\n", contact1.ID)
	}

	if err == nil {
		t.Errorf("Expected Error %s, but got <nil>\n", "Contact Already Exists")
	}
}

func TestOrderOfVisits(t *testing.T) {
	expected := []int{159, 158, 155, 147, 146, 145, 144, 142, 141, 140, 138, 132, 131, 129, 128, 127, 126, 120, 119, 117, 116, 113, 112, 111, 108, 106, 105, 104, 103, 102, 98, 96, 94, 91, 89, 86, 84, 82, 79, 74, 72, 71, 67, 66, 65, 62, 59, 58, 57, 53, 51, 49, 47, 43, 38, 31, 29, 28, 25, 23, 21, 20, 19, 18, 14, 13, 12, 10, 9, 6, 5, 3, 1, 0, 2, 4, 7, 8, 11, 15, 16, 17, 22, 24, 26, 27, 30, 32, 33, 34, 35, 36, 37, 39, 40, 41, 42, 44, 45, 46, 48, 50, 52, 54, 55, 56, 60, 61, 63, 64, 68, 69, 70, 73, 75, 76, 77, 78, 80, 81, 83, 85, 87, 88, 90, 92, 93, 95, 97, 99, 100, 101, 107, 109, 110, 114, 115, 118, 121, 122, 123, 124, 125, 130, 133, 134, 135, 136, 137, 139, 143, 148, 149, 150, 151, 152, 153, 154, 156, 157}
	id, _ := From("C80F741BC1B397C54A54858E4E2A8840B2BC766B")
	routing := NewRoutingTable(nil)

	out := routing.determineOrderOfVisits(id)

	if len(expected) != len(out) {
		t.Errorf("Expected %d but got %d\n", len(expected), len(out))
	}

	for i, x := range expected {
		if out[i] != x {
			t.Errorf("Expected %d at index %d, but got %d\n", x, i, out[i])
		}
	}

}

func TestRoutingGet3closestContacts(t *testing.T) {
	id, _ := From("395754ecb968b3d40ab6ea17322edd4b84012938")
	lookupId, _ := From("16bcc112cd86800edfd11b0f7d2a2c476bd34f22")
	contactIds := []string{
		"8f2d6ae2378dda228d3bd39c41a4b6f6f538a41a", // xor to lookup: (2) 1001100110010001101010111111000011111010000010110101101000101100010100101110101011001000100100110011110010001110100110101011000110011110111010111110101100111000
		"28f787e3b60f99fb29b14266c40b536d6037307e", // xor to lookup: (0) 11111001001011010001101111000101111011100010010001100111110101111101100110000001011001011010011011100100100001011111110010101000001011111001000111111101011100
		"b4945c02ddd3d4484ed7200107b46f65f5300305", // xor to lookup: (3) 1010001000101000100111010001000000010000010101010101010001000110100100010000011000111011000011100111101010011110010000110010001010011110111000110100110000100111
		"dc03f8f281c7118225901c8655f788cd84e3f449", // xor to lookup: (4) 1100101010111111001110011110000001001100010000011001000110001100111110100100000100000111100010010010100011011101101001001000101011101111001100001011101101101011
		"9d079f19f9edca7f8b2f5ce58624b55ffec2c4f3", // xor to lookup: (1) 1000101110111011010111100000101100110100011010110100101001110001010101001111111001000111111010101111101100001110100110010001100010010101000100011000101111010001
	}

	routing := NewRoutingTable(id)

	for _, hex := range contactIds {
		c := generateContactFrom(hex)
		_, index, _ := routing.Add(c)
		fmt.Printf("Added at index %d\n", index)
	}

	c := routing.getXClosestContacts(3, lookupId)

	if len(c) != 3 {
		t.Errorf("Expected length of %d, but got %d\n", 3, len(c))
	}

	if c[0].ID.String() != "28f787e3b60f99fb29b14266c40b536d6037307e" {
		t.Errorf("Expected at index (%d) %s, but got %s\n", 0, "28f787e3b60f99fb29b14266c40b536d6037307e", c[0].ID.String())
	}

	if c[1].ID.String() != "9d079f19f9edca7f8b2f5ce58624b55ffec2c4f3" {
		t.Errorf("Expected at index (%d) %s, but got %s\n", 1, "9d079f19f9edca7f8b2f5ce58624b55ffec2c4f3", c[1].ID.String())

	}

	if c[2].ID.String() != "8f2d6ae2378dda228d3bd39c41a4b6f6f538a41a" {
		t.Errorf("Expected at index (%d) %s, but got %s\n", 2, "8f2d6ae2378dda228d3bd39c41a4b6f6f538a41a", c[2].ID.String())

	}
}

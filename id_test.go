package gokad

import (
	"fmt"
	"testing"
)

func TestGenerateID(t *testing.T) {
	id := GenerateRandomID()
	hex := fmt.Sprintf("%s", id)
	t.Logf("ID %s\n", hex)
	id2, _ := From(hex)

	if !id.Equal(id2) {
		t.Errorf("Expected id (%s) to equal id2 (%s)\n", id, id2)
	}
}

func TestStringID(t *testing.T) {
	id, _ := From("0c204d39600fddd3f1f20ca8007e91c7d0293b1c")

	if id.String() != "0c204d39600fddd3f1f20ca8007e91c7d0293b1c" {
		t.Errorf("Expected to be %s, but got %s", "0c204d39600fddd3f1f20ca8007e91c7d0293b1c", id.String())
	}
}

func TestCompareDistanceTo(t *testing.T) {
	id1, e := From("0c204d39600fddd3f1f20ca8007e91c7d0293b1c")
	id2, e2 := From("0c204d39600fddd3f1f20ca8007e91c7d0293b1f") // distance is 11
	id3, e3 := From("0c204d39600fddd3f1f20fa8007e91c7d0293b1c") // distance is much larger

	if e != nil || e2 != nil || e3 != nil {
		t.Errorf("Error in constructing Id %s %s %s", e, e2, e3)
	}

	res := id1.CompareDistanceTo(id2, id3)

	if res < 1 {
		t.Errorf("Expected res to be 1, but got %d", res)
	}
}

func TestGetBitAt(t *testing.T) {
	// 1100100000001111011101000001101111000001101100111001011111000101010010100101010010000101100011100100111000101010100010000100000010110010101111000111011001101011
	id, _ := From("c80f741bc1b397c54a54858e4e2a8840b2bc766b")

	cases := []struct {
		In  uint
		Out int
	}{
		{
			0,
			1,
		},
		{
			2,
			0,
		},
		{
			3,
			0,
		},
		{
			160,
			1,
		},
		{
			159,
			1,
		},
	}

	for _, c := range cases {
		bit := id.GetBitAt(c.In)

		if bit != c.Out {
			t.Errorf("Expected bit to be %d, but got %d [Input: %d]", c.Out, bit, c.In)
		}
	}

}

package gokad

import "testing"

func TestBucketIndex(t *testing.T) {
	cases := []struct {
		IN  int
		OUT int
	}{
		{
			IN:  0,
			OUT: 159,
		},
		{
			IN:  159,
			OUT: 0,
		},
		{
			IN:  10,
			OUT: 149,
		},
	}

	for _, c := range cases {
		index := bucketIndex(c.IN)
		if index != c.OUT {
			t.Errorf("Expected %d but got %d\n", c.OUT, index)
		}
	}
}

func TestSortDistance(t *testing.T) {
	root, _ := From("C80F741BC1B397C54A54858E4E2A8840B2BC766B")
	id1, _ := From("D80F741BC1B397C54A54858E4E2A8840B2BC766B")
	id2, _ := From("F70F741BC1B397C54A54858E4E2A8840B2BC766B")
	id3, _ := From("A70F441BC1B397C54A54858E4E2A8840B2BC766B")

	delta1 := root.DistanceTo(id1)
	delta2 := root.DistanceTo(id2)
	delta3 := root.DistanceTo(id3)

	input := []*Distance{delta3, delta1, delta2}
	sort(input)

	for i, d := range input {
		if i == 0 {
			if d != delta1 {
				t.Errorf("Expected (delta1) %s, but got %s\n", delta1, d)
			}
			continue
		}

		if i == 1 {
			if d != delta2 {
				t.Errorf("Expected (delta2) %s, but got %s\n", delta2, d)
			}
			continue
		}

		if i == 2 {
			if d != delta3 {
				t.Errorf("Expected (delta3) %s, but got %s\n", delta3, d)
			}
			continue
		}
	}

}

package gokad

// bucketIndex is a utility function to get the bucket index
// there are 160 buckets in a routing table.
// if i have a distance of 01001111 in a 8 bit address space.
// I often read/loop through distances starting at index 0, but in reality that is bucket 159
// so in this case,
// Bit:    0  1  0  0  1  1  1  1
// Bucket: 7  6  5  4  3  2  1  0
func bucketIndex(index int) int {
	return (MaxRoutingTableSize - 1) - index
}

// compareDistance compares 2 distances to each other
// return 1 if d1 is larger, -1 if d2 is larger and 0 if they are the same
func compareDistance(d1, d2 *Distance) int {
	for i := 0; i < MaxCapacity; i++ {
		if d1.buf[i] > d2.buf[i] {
			return 1
		}

		if d1.buf[i] < d2.buf[i] {
			return -1
		}
	}

	return 0
}

func sort(x []*Distance) {
	for i := 0; i < len(x); i++ {
		for j := i + 1; j < len(x); j++ {
			if compareDistance(x[i], x[j]) > 0 {
				smaller := x[j]
				x[j] = x[i]
				x[i] = smaller
			}
		}
	}

}

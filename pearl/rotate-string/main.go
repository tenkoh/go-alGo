package main

// rotate rotates a string by shifts. this algorithm save memory, no allocation.
func rotate(b []byte, shift int) {
	// split b into two parts
	// b[:shift], b[shift:]
	// if shift is bigger than len(b), shift = shift % len(b)
	// if shift is negative, shift = len(b) - shift
	if len(b) <= 1 {
		return
	}

	for shift < 0 {
		shift += len(b)
	}
	shift = shift % len(b)

	// rotate(b) = r(r(left)+r(right))
	reverse(b, 0, shift-1)
	reverse(b, shift, len(b)-1)
	reverse(b, 0, len(b)-1)
	return
}

func reverse(b []byte, start, last int) {
	if len(b) == 0 {
		return
	}
	if start >= last {
		return
	}
	if start < 0 || last >= len(b) {
		return
	}

	// replacement range
	// [0 1] -> 1回実施 (last-start) = 1
	// [0 1 2] -> 1回実施 (last-start) = 2
	// [0 1 2 3] -> 2回実施 (last-start) = 3
	// ということでw := last - start + 1 として w/2で回数を決める
	w := last - start + 1
	for i := 0; i < w/2; i++ {
		b[start+i], b[last-i] = b[last-i], b[start+i]
	}
	return
}

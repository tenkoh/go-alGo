package main

// rotate rotates a string by shifts. this algorithm save memory, no allocation.
func rotate(s string, shift int) string {
	return ""
}

func reverse(s string, start, last int) string {
	if last > len(s)-1 || start > last {
		panic("invalid index")
	}
	if start == last {
		return s
	}

	w := last - start
	for i := 0; i < w/2; i++ {

	}
	return s
}

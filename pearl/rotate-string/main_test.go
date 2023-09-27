package main

import "testing"

func TestRotate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		shift int
		want  string
	}{
		{"empty string", "", 1, ""},
		{"single character", "a", 1, "a"},
		{"single character", "a", 2, "a"},
		{"some characters", "abc", 1, "bca"},
		{"some characters", "abc", 2, "cab"},
		{"some characters", "abc", 3, "abc"},
		{"some characters", "abc", 4, "bca"},
		{"minus shift", "abc", -1, "cab"},
		{"minus shift", "abc", -2, "bca"},
		{"minus shift", "abc", -3, "abc"},
		{"minus shift", "abc", -4, "cab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rotate(tt.input, tt.shift)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

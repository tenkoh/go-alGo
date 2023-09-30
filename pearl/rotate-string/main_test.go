package main

import (
	"slices"
	"testing"
)

func TestRotate(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		shift int
		want  []byte
	}{
		{"empty string", []byte(""), 0, []byte("")},
		{"single character", []byte("a"), 0, []byte("a")},
		{"single character", []byte("a"), 1, []byte("a")},
		{"single character", []byte("a"), 2, []byte("a")},
		{"some characters, no rotate", []byte("abc"), 0, []byte("abc")},
		{"some characters, +1", []byte("abc"), 1, []byte("bca")},
		{"some characters, +2", []byte("abc"), 2, []byte("cab")},
		{"some characters, +3", []byte("abc"), 3, []byte("abc")},
		{"some characters, +4=+1", []byte("abc"), 4, []byte("bca")},
		{"some characters, -1=+2", []byte("abc"), -1, []byte("cab")},
		{"some characters, -2=+1", []byte("abc"), -2, []byte("bca")},
		{"some characters, -3=0", []byte("abc"), -3, []byte("abc")},
		{"some characters, -4=+2", []byte("abc"), -4, []byte("cab")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rotate(tt.input, tt.shift)
			if !slices.Equal(tt.input, tt.want) {
				t.Errorf("got %q, want %q", tt.input, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		start int
		last  int
		want  []byte
	}{
		{"empty string", []byte(""), 0, 0, []byte("")},
		{"bad start", []byte("abc"), -1, 0, []byte("abc")},
		{"bad last", []byte("abc"), 0, 3, []byte("abc")},
		{"bad start and last", []byte("abc"), 2, 1, []byte("abc")},
		{"single character", []byte("a"), 0, 0, []byte("a")},
		{"some characters", []byte("abc"), 0, 0, []byte("abc")},
		{"some characters", []byte("abc"), 0, 1, []byte("bac")},
		{"some characters", []byte("abc"), 0, 2, []byte("cba")},
		{"some characters", []byte("abc"), 1, 1, []byte("abc")},
		{"some characters", []byte("abc"), 1, 2, []byte("acb")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reverse(tt.input, tt.start, tt.last)
			if !slices.Equal(tt.input, tt.want) {
				t.Errorf("got %q, want %q", tt.input, tt.want)
			}
		})
	}
}

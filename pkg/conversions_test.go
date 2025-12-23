package pkg

import (
	"testing"
)

func TestSIToBytes(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{"1.5 KB", 1536},
		{"1.0 MB", 1048576},
		{"1.0 GB", 1073741824},
		{"1.0 TB", 1099511627776},
		{"1.0 PB", 1125899906842624},
		{"1.0 EB", 1152921504606846976},
		{"0.0 B", 0},
		{"0 B", 0},
		{"1023 B", 1023},
	}
	for _, tt := range tests {
		got, err := SIToBytes(tt.in)
		if err != nil {
			t.Fatalf("SIToBytes(%q) error: %v", tt.in, err)
		}
		if got != tt.want {
			t.Fatalf("SIToBytes(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestBytesToSI(t *testing.T) {
	tests := []struct {
		in   float64
		want string
	}{
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
		{1125899906842624, "1.0 PB"},
		{1152921504606846976, "1.0 EB"},
		{0, "0.0 B"},
		{1023, "1023.0 B"},
	}
	for _, tt := range tests {
		got := BytesToSI(tt.in)
		if got != tt.want {
			t.Fatalf("BytesToSI(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestSecondsToText(t *testing.T) {
	tests := []struct {
		in   int64
		want string
	}{
		{3601, "1 hour, 1 second"},
		{2443044, "28 days, 6 hours, 37 minutes, 24 seconds"},
		{350735, "4 days, 1 hour, 25 minutes, 35 seconds"},
		{5114048, "1 month, 29 days, 4 hours, 34 minutes, 8 seconds"},
	}
	for _, tt := range tests {
		got := SecondsToText(tt.in)
		if got != tt.want {
			t.Fatalf("SecondsToText(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

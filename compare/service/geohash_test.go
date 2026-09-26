package service

import "testing"

func TestEncodeGeohashKnownPoint(t *testing.T) {
	// UNBF campus should be stable at precision 5.
	h := EncodeGeohash(45.9458, -66.6414, 5)
	if len(h) != 5 {
		t.Fatalf("expected length 5, got %q", h)
	}
	for _, c := range h {
		if !geohashCharValid(byte(c)) {
			t.Fatalf("invalid char %c in %q", c, h)
		}
	}
}

func geohashCharValid(c byte) bool {
	for i := 0; i < len(geohashBase32); i++ {
		if geohashBase32[i] == c {
			return true
		}
	}
	return false
}

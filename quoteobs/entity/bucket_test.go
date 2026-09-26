package entity

import "testing"

func TestBasketSubtotalBucket(t *testing.T) {
	if BasketSubtotalBucket(1350) != 1000 {
		t.Fatalf("expected bucket 1000 for 1350")
	}
	if BasketSubtotalBucket(500) != 500 {
		t.Fatalf("expected 500")
	}
}

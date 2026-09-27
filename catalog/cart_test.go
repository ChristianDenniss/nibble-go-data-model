package catalog

import (
	"errors"
	"math"
	"testing"
)

func TestWholeCartComparison(t *testing.T) {
	c := Build(fixture())
	r := c.Restaurants[0]
	req := CartRequest{RestaurantID: r.ID, Lines: []CartLine{{ItemID: r.Items[0].ID, Quantity: 2}, {ItemID: r.Items[1].ID, Quantity: 3}}}
	got, err := CompareCart(c, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Providers) != 2 {
		t.Fatal("missing provider")
	}
	for _, p := range got.Providers {
		if len(p.Lines) != 2 || p.Total != nil || p.Delivery != nil || p.Service != nil || p.Tax != nil {
			t.Fatal("lost lines or fabricated checkout fees")
		}
		if p.Provider == "DoorDash" {
			if !p.Complete || *p.Subtotal != 6110 || !p.StartingPrice {
				t.Fatal("incorrect complete basket", p)
			}
		}
		if p.Provider == "Uber Eats" {
			if p.Complete || p.Subtotal != nil || p.KnownSubtotal != 3060 || p.Lines[1].LineTotal != nil {
				t.Fatal("missing item treated as free", p)
			}
		}
	}
	req.Lines = req.Lines[:1]
	got, err = CompareCart(c, req)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range got.Providers {
		if !p.Complete || p.Subtotal == nil {
			t.Fatal("complete cart not priced")
		}
	}
}
func TestCartValidationAndOverflow(t *testing.T) {
	c := Build(fixture())
	r := c.Restaurants[0]
	line := CartLine{ItemID: r.Items[0].ID, Quantity: 1}
	for _, req := range []CartRequest{{RestaurantID: r.ID}, {RestaurantID: "other", Lines: []CartLine{line}}, {RestaurantID: r.ID, Lines: []CartLine{{ItemID: line.ItemID, Quantity: 0}}}, {RestaurantID: r.ID, Lines: []CartLine{{ItemID: line.ItemID, Quantity: 51}}}, {RestaurantID: r.ID, Lines: []CartLine{line, line}}, {RestaurantID: r.ID, Lines: []CartLine{{ItemID: "missing", Quantity: 1}}}} {
		if _, err := CompareCart(c, req); !errors.Is(err, ErrInvalidCart) {
			t.Fatal("invalid cart accepted", req)
		}
	}
	r.Items[0].Offers[0].Amount = price(math.MaxInt64)
	if _, err := CompareCart(c, CartRequest{RestaurantID: r.ID, Lines: []CartLine{{ItemID: line.ItemID, Quantity: 2}}}); !errors.Is(err, ErrInvalidCart) {
		t.Fatal("overflow accepted")
	}
}

func TestRankDoesNotCallIncompleteOrStartingBasketsCheapest(t *testing.T) {
	providers := []ProviderCart{
		{Provider: "Incomplete", KnownSubtotal: 1},
		{Provider: "Starting", Complete: true, StartingPrice: true, Subtotal: price(2)},
		{Provider: "DoorDash", Complete: true, Subtotal: price(1300)},
		{Provider: "SkipTheDishes", Complete: true, Subtotal: price(1200)},
		{Provider: "Uber Eats", Complete: true, Subtotal: price(1200)},
	}
	RankCarts(providers)
	if providers[0].Provider != "SkipTheDishes" || !providers[0].LowestListedSubtotal || !providers[1].LowestListedSubtotal {
		t.Fatal("ties not ranked correctly")
	}
	for _, p := range providers {
		if (!p.Complete || p.StartingPrice) && p.LowestListedSubtotal {
			t.Fatal("unknown basket ranked cheapest")
		}
		if p.Total != nil {
			t.Fatal("delivered total invented")
		}
	}
}

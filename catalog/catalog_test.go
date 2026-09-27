package catalog

import (
	"encoding/json"
	"strings"
	"testing"
)

func price(n int64) *int64 { return &n }
func fixture() Bundle {
	uber := Store{ID: "ss_ubereats_0a5ba1d349096949787465cf", Name: "Taco Boyz", URL: "https://www.ubereats.com/ca/store/taco-boyz/gnecSBVbV9ygUE84RWhwVA", Address: "520 Smythe St #2A", SourceAge: "3 weeks", Items: []MenuItem{{Name: "Large Nachos", Amount: price(1530), Currency: "CAD"}}}
	dd := Store{ID: "ss_doordash_23391251", Name: "Taco Boyz", URL: "https://www.doordash.com/store/taco-boyz-fredericton-23391251/", ObservedAt: "2026-09-26T12:00:00Z", Items: []MenuItem{{Name: "Large Nachos™ [570.0 Cals]", Amount: price(1555), Currency: "CAD", PriceKind: "from"}, {Name: "Small Nachos", Amount: price(1000), Currency: "CAD"}}}
	return Bundle{Version: 1, Providers: map[string]Snapshot{"Uber Eats": {City: "Fredericton", Region: "NB", RetrievedAt: "2026-09-26T12:00:00Z", Stores: []Store{uber}}, "DoorDash": {City: "Fredericton", Region: "NB", RetrievedAt: "2026-09-26T12:00:00Z", Stores: []Store{dd}}}}
}
func TestReviewedBranchesAndPrices(t *testing.T) {
	b := fixture()
	if err := b.Validate(); err != nil {
		t.Fatal(err)
	}
	got := Build(b)
	if len(got.Restaurants) != 1 {
		t.Fatalf("branches not merged: %d", len(got.Restaurants))
	}
	r := got.Restaurants[0]
	if len(r.Items) != 2 || len(r.Items[0].Offers) != 2 {
		t.Fatal("items not resolved")
	}
	offers := r.Items[0].Offers
	if *offers[0].Amount != 1555 || *offers[1].Amount != 1530 || offers[0].PriceKind != "from" {
		t.Fatal("price changed")
	}
	missing := r.Items[1].Offers[1]
	if missing.Provider != "Uber Eats" || missing.Amount != nil || missing.Status != "unavailable" || missing.Delivery != nil {
		t.Fatal("missing price fabricated")
	}
	raw, _ := json.Marshal(got)
	for _, private := range []string{"observedAt", "retrievedAt", "sourceAge", "directoryObservedAt", "3 weeks", "snapshot"} {
		if strings.Contains(string(raw), private) {
			t.Fatalf("public provenance leaked: %s", private)
		}
	}
	dd := b.Providers["DoorDash"]
	dd.Stores[0].ID = "ss_doordash_999"
	dd.Stores[0].URL = "https://www.doordash.com/store/taco-boyz-fredericton-999/"
	b.Providers["DoorDash"] = dd
	if len(Build(b).Restaurants) != 2 {
		t.Fatal("unreviewed branch merged by name")
	}
}
func TestItemIdentityPreservesVariants(t *testing.T) {
	for _, pair := range [][2]string{{"Big Mac", "Big Mac Extra Value Meal"}, {"12 inch Pizza", "14 inch Pizza"}, {"Chicken on the Rocks", "On the Rocks"}, {"10 McNuggets", "20 McNuggets"}} {
		if ItemKey(pair[0]) == ItemKey(pair[1]) {
			t.Fatal("variant merged", pair)
		}
	}
	if ItemKey("Chef’s Bowl™ [123.0 Cals]") != ItemKey("Chef's Bowl") {
		t.Fatal("cosmetic difference not normalized")
	}
	b := fixture()
	u := b.Providers["Uber Eats"]
	u.Stores[0].Items = append(u.Stores[0].Items, u.Stores[0].Items[0])
	b.Providers["Uber Eats"] = u
	for _, it := range Build(b).Restaurants[0].Items {
		count := 0
		for _, o := range it.Offers {
			if o.Amount != nil {
				count++
			}
		}
		if count > 1 {
			t.Fatal("ambiguous price merged")
		}
	}
}
func TestValidation(t *testing.T) {
	for _, mutate := range []func(*Store){func(s *Store) { s.Items[0].Amount = price(-1) }, func(s *Store) { s.Items[0].Currency = "USD" }, func(s *Store) { s.URL = "https://evil.example/menu" }, func(s *Store) { s.SourceAge = "" }, func(s *Store) { s.ObservedAt = "yesterday" }} {
		b := fixture()
		s := b.Providers["Uber Eats"]
		mutate(&s.Stores[0])
		b.Providers["Uber Eats"] = s
		if b.Validate() == nil {
			t.Fatal("invalid input accepted")
		}
	}
}

func TestSkipJoinsOnlyReviewedBranch(t *testing.T) {
	b := fixture()
	url := "https://www.skipthedishes.com/taco-boyz-520"
	b.Providers["SkipTheDishes"] = Snapshot{City: "Fredericton", Region: "NB", RetrievedAt: "2026-09-27T00:00:00Z", Stores: []Store{{ID: StoreID("SkipTheDishes", url), Name: "Taco Boyz", URL: url, SourceAge: "2 weeks", Items: []MenuItem{{Name: "Large Nachos", Amount: price(1500), Currency: "CAD"}}}}}
	if err := b.Validate(); err != nil {
		t.Fatal(err)
	}
	c := Build(b)
	if len(c.Restaurants) != 1 || len(c.Restaurants[0].Providers) != 3 || len(c.Restaurants[0].Items[0].Offers) != 3 {
		t.Fatal("Skip branch not merged")
	}
}

func TestExpandedMcDonaldsBranchesRemainSeparate(t *testing.T) {
	prospect := "ss_ubereats_b516b917b6ed6c6e5c3ba08c"
	main := "ss_ubereats_49846df7683e641f6ac1cc0c"
	item := func(n int64) []MenuItem { return []MenuItem{{Name: "Big Mac", Amount: price(n), Currency: "CAD"}} }
	b := Bundle{Version: 1, Providers: map[string]Snapshot{
		"Uber Eats":     {Stores: []Store{{ID: prospect, Name: "McDonald's Prospect", Items: item(839)}, {ID: main, Name: "McDonald's Main", Items: item(900)}}},
		"DoorDash":      {Stores: []Store{{ID: "ss_doordash_1139998", Name: "McDonald's", Items: item(859)}}},
		"SkipTheDishes": {Stores: []Store{{ID: "ss_skip_e9300610bea06b1a5b53937a", Name: "McDonald's", Items: item(879)}, {ID: "ss_skip_150651ddd49a40b4bfd9fbfa", Name: "McDonald's", Items: item(920)}}},
	}}
	got := Build(b)
	if len(got.Restaurants) != 2 {
		t.Fatalf("expected two distinct branches, got %d", len(got.Restaurants))
	}
	for _, restaurant := range got.Restaurants {
		want := 2
		if restaurant.ID == "catalog-"+prospect {
			want = 3
		}
		if len(restaurant.Items) != 1 || len(restaurant.Items[0].Offers) != want {
			t.Fatalf("wrong provider association for %s", restaurant.ID)
		}
	}
}

func TestReviewedMappingsHaveEvidenceAndNoAliasChains(t *testing.T) {
	var mappings map[string]struct {
		Canonical   string `json:"canonicalStoreId"`
		EvidenceURL string `json:"evidenceUrl"`
		Evidence    string `json:"evidence"`
	}
	raw, _ := rules.ReadFile("branch_matches.json")
	if err := json.Unmarshal(raw, &mappings); err != nil {
		t.Fatal(err)
	}
	for source, mapping := range mappings {
		if mapping.Canonical == source || mapping.EvidenceURL == "" || mapping.Evidence == "" {
			t.Fatalf("invalid reviewed mapping: %s", source)
		}
		if _, chain := mappings[mapping.Canonical]; chain {
			t.Fatalf("alias chain must resolve to a root: %s", source)
		}
	}
	if _, merged := mappings["ss_doordash_2518559"]; merged {
		t.Fatal("Subway conflicting unit addresses must not be merged")
	}
}

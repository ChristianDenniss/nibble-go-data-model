package service

import (
	"context"
	"testing"
	"time"

	"github.com/ChristianDenniss/go-data-model/serviceability/entity"
)

type testStatusRepo struct{ status entity.StoreStatus }

func (r testStatusRepo) Get(context.Context, string) (entity.StoreStatus, error) {
	if r.status.SourceStoreID == "" {
		return entity.StoreStatus{}, entity.ErrNotFound
	}
	return r.status, nil
}
func (testStatusRepo) Upsert(context.Context, entity.StoreStatus) error { return nil }

type testAreaRepo struct{ areas []entity.ServiceArea }

func (r testAreaRepo) ListForPath(context.Context, string, string, string) ([]entity.ServiceArea, error) {
	return r.areas, nil
}
func (testAreaRepo) Upsert(context.Context, entity.ServiceArea) error { return nil }

func TestDeliverableAtUsesPolygonAndHoles(t *testing.T) {
	service := New(testStatusRepo{}, testAreaRepo{areas: []entity.ServiceArea{{Geometry: `{"type":"Polygon","coordinates":[[[-64,45],[-63,45],[-63,46],[-64,46],[-64,45]]]}`}}})

	covered, reason, err := service.DeliverableAt(context.Background(), "store", "delivery", "", "f80", 45.1, -63.5)
	if err != nil || !covered || reason != "" {
		t.Fatalf("expected polygon coverage, got covered=%v reason=%q err=%v", covered, reason, err)
	}

	holeGeometry := `{"type":"Polygon","coordinates":[[[-64,45],[-63,45],[-63,46],[-64,46],[-64,45]],[[-63.8,45.2],[-63.2,45.2],[-63.2,45.8],[-63.8,45.8],[-63.8,45.2]]]}`
	if geometryCovered(holeGeometry, "f80", 45.5, -63.5) {
		t.Fatal("expected polygon hole to be excluded")
	}
}

func TestDeliverableAtRejectsPausedAndStaleStores(t *testing.T) {
	paused := New(testStatusRepo{status: entity.StoreStatus{SourceStoreID: "store", OpenNow: true, Paused: true}}, testAreaRepo{})
	covered, reason, err := paused.DeliverableAt(context.Background(), "store", "delivery", "", "f80", 45, -63)
	if err != nil || covered || reason != "store_not_available" {
		t.Fatalf("expected paused store to be unavailable, got covered=%v reason=%q err=%v", covered, reason, err)
	}

	stale := New(testStatusRepo{status: entity.StoreStatus{SourceStoreID: "store", OpenNow: true, ObservedAt: time.Now().Add(-time.Hour)}}, testAreaRepo{})
	covered, reason, err = stale.DeliverableAt(context.Background(), "store", "delivery", "", "f80", 45, -63)
	if err != nil || covered || reason != "status_stale" {
		t.Fatalf("expected stale store to be unknown, got covered=%v reason=%q err=%v", covered, reason, err)
	}
}

func TestDeliverableAtReportsUnknownCoverageWithoutAreas(t *testing.T) {
	service := New(testStatusRepo{}, testAreaRepo{})
	covered, reason, err := service.DeliverableAt(context.Background(), "store", "delivery", "", "f80", 45, -63)
	if err != nil || !covered || reason != "coverage_unknown" {
		t.Fatalf("expected unknown coverage to remain visible, got covered=%v reason=%q err=%v", covered, reason, err)
	}
}

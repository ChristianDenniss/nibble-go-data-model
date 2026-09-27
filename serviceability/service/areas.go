package service

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"
	"time"

	fulfillmententity "github.com/ChristianDenniss/go-data-model/fulfillment/entity"
	serviceabilityentity "github.com/ChristianDenniss/go-data-model/serviceability/entity"
)

const DefaultStatusMaxAge = 30 * time.Minute

// Deliverable checks store status and optional service_areas geohash prefixes (comma-separated).
// geometry "*" or empty with a row present means entire area; no rows means unknown (allow).
func (s *Service) Deliverable(ctx context.Context, sourceStoreID, fulfillment, deliveryExecutor, dropoffGeohash string) (bool, string, error) {
	return s.DeliverableAt(ctx, sourceStoreID, fulfillment, deliveryExecutor, dropoffGeohash, 0, 0)
}

func (s *Service) DeliverableAt(ctx context.Context, sourceStoreID, fulfillment, deliveryExecutor, dropoffGeohash string, latitude, longitude float64) (bool, string, error) {
	if st, err := s.status.Get(ctx, sourceStoreID); err == nil {
		if !st.ObservedAt.IsZero() && time.Since(st.ObservedAt) > DefaultStatusMaxAge {
			return false, "status_stale", nil
		}
		if st.Paused || !st.OpenNow {
			return false, "store_not_available", nil
		}
	} else if !errors.Is(err, serviceabilityentity.ErrNotFound) {
		return false, "", err
	}

	if fulfillment != fulfillmententity.ModeDelivery {
		return true, "", nil
	}
	if dropoffGeohash == "" {
		return false, "dropoff_required", nil
	}

	areas, err := s.areas.ListForPath(ctx, sourceStoreID, fulfillment, deliveryExecutor)
	if err != nil {
		return false, "", err
	}
	if len(areas) == 0 {
		return true, "coverage_unknown", nil
	}
	for _, a := range areas {
		if geometryCovered(a.Geometry, dropoffGeohash, latitude, longitude) {
			return true, "", nil
		}
	}
	return false, "outside_service_area", nil
}

func geometryCovered(geometry, geohash string, latitude, longitude float64) bool {
	geometry = strings.TrimSpace(geometry)
	if strings.HasPrefix(geometry, "{") || strings.HasPrefix(geometry, "[") {
		return geoJSONContains(geometry, latitude, longitude)
	}
	return geohashCovered(geometry, geohash)
}

type geoJSONGeometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

func geoJSONContains(raw string, latitude, longitude float64) bool {
	var geometry geoJSONGeometry
	if json.Unmarshal([]byte(raw), &geometry) != nil || geometry.Type != "Polygon" {
		return false
	}
	var rings [][][]float64
	if json.Unmarshal(geometry.Coordinates, &rings) != nil || len(rings) == 0 {
		return false
	}
	return pointInRing(rings[0], longitude, latitude) && !insideHole(rings[1:], longitude, latitude)
}

func insideHole(holes [][][]float64, longitude, latitude float64) bool {
	for _, hole := range holes {
		if pointInRing(hole, longitude, latitude) {
			return true
		}
	}
	return false
}

func pointInRing(ring [][]float64, x, y float64) bool {
	inside := false
	for i, j := 0, len(ring)-1; i < len(ring); j, i = i, i+1 {
		if len(ring[i]) < 2 || len(ring[j]) < 2 {
			continue
		}
		xi, yi := ring[i][0], ring[i][1]
		xj, yj := ring[j][0], ring[j][1]
		intersects := (yi > y) != (yj > y) && x < (xj-xi)*(y-yi)/(yj-yi+math.SmallestNonzeroFloat64)+xi
		if intersects {
			inside = !inside
		}
	}
	return inside
}

func geohashCovered(geometry, geohash string) bool {
	geometry = strings.TrimSpace(geometry)
	if geometry == "" || geometry == "*" {
		return true
	}
	for _, p := range strings.Split(geometry, ",") {
		p = strings.TrimSpace(p)
		if p != "" && strings.HasPrefix(geohash, p) {
			return true
		}
	}
	return false
}

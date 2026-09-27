package service

import (
	"math"
	"sort"

	"github.com/ChristianDenniss/go-data-model/home/entity"
	restaurantentity "github.com/ChristianDenniss/go-data-model/restaurant/entity"
	storefrontentity "github.com/ChristianDenniss/go-data-model/storefront/entity"
)

// Popular ranks by rating average weighted by how many people rated (log-scaled).
// Takes only the catalog: organic rails never see merchandising (D28).
func Popular(catalog storefrontentity.Catalog, limit int) []entity.FeedItem {
	ranked := append([]restaurantentity.Restaurant(nil), catalog.Restaurants...)
	sort.SliceStable(ranked, func(i, j int) bool {
		a, b := popularity(ranked[i].Rating), popularity(ranked[j].Rating)
		if a != b {
			return a > b
		}
		return ranked[i].ID < ranked[j].ID
	})
	return toItems(ranked, limit, func(restaurantentity.Restaurant) string { return "" })
}

func popularity(r restaurantentity.Rating) float64 {
	return r.Average * math.Log1p(float64(r.Count))
}

// Recommended ranks by overlap with cuisines and categories from the account's past orders,
// then by rating. With no order history it falls back to highest rated.
func Recommended(catalog storefrontentity.Catalog, limit int) []entity.FeedItem {
	byID := map[string]restaurantentity.Restaurant{}
	for _, r := range catalog.Restaurants {
		byID[r.ID] = r
	}
	cuisineWeight := map[string]int{}
	categoryWeight := map[string]int{}
	for _, o := range catalog.Orders {
		r, ok := byID[o.RestaurantID]
		if !ok {
			continue
		}
		for _, id := range r.CuisineIDs {
			cuisineWeight[id]++
		}
		for _, id := range r.CategoryIDs {
			categoryWeight[id]++
		}
	}
	cuisineName := map[string]string{}
	for _, c := range catalog.Cuisines {
		cuisineName[c.ID] = c.Name
	}

	score := func(r restaurantentity.Restaurant) float64 {
		s := r.Rating.Average
		for _, id := range r.CuisineIDs {
			s += 2 * float64(cuisineWeight[id])
		}
		for _, id := range r.CategoryIDs {
			s += float64(categoryWeight[id])
		}
		return s
	}

	ranked := append([]restaurantentity.Restaurant(nil), catalog.Restaurants...)
	sort.SliceStable(ranked, func(i, j int) bool {
		a, b := score(ranked[i]), score(ranked[j])
		if a != b {
			return a > b
		}
		if ranked[i].Rating.Count != ranked[j].Rating.Count {
			return ranked[i].Rating.Count > ranked[j].Rating.Count
		}
		return ranked[i].ID < ranked[j].ID
	})

	return toItems(ranked, limit, func(r restaurantentity.Restaurant) string {
		best, bestWeight := "", 0
		for _, id := range r.CuisineIDs {
			if cuisineWeight[id] > bestWeight && cuisineName[id] != "" {
				best, bestWeight = cuisineName[id], cuisineWeight[id]
			}
		}
		if best != "" {
			return "Because you like " + best
		}
		if len(catalog.Orders) == 0 && r.Rating.Average >= 4.5 {
			return "Highly rated"
		}
		return ""
	})
}

func toItems(ranked []restaurantentity.Restaurant, limit int, reason func(restaurantentity.Restaurant) string) []entity.FeedItem {
	if len(ranked) > limit {
		ranked = ranked[:limit]
	}
	out := make([]entity.FeedItem, len(ranked))
	for i, r := range ranked {
		out[i] = entity.FeedItem{RestaurantID: r.ID, Reason: reason(r)}
	}
	return out
}

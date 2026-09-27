package entity

import location "github.com/ChristianDenniss/go-data-model/location/entity"

type Restaurant struct {
	ID          string
	Name        string
	ImageURL    string
	Location    location.Location
	CuisineIDs  []string
	CategoryIDs []string
	Rating      Rating
	// Phone takes call-in orders when set. Empty means the restaurant does not take phone orders.
	Phone string
	// AppURL is the restaurant's own ordering app (store listing or deep link). Empty means no app.
	AppURL string
	// Hours holds both the store and delivery schedules. Empty means hours are unknown.
	Hours []Hours
}

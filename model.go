package model

import "time"

type Money struct {
	AmountCents int64
	Currency    string
}

type Location struct {
	Latitude  float64
	Longitude float64
	Address   string
}

type Provider struct {
	ID   string
	Name string
}

type Restaurant struct {
	ID       string
	Name     string
	Location Location
}

type MenuItem struct {
	ID           string
	RestaurantID string
	Name         string
}

type Offer struct {
	ID           string
	RestaurantID string
	ProviderID   string
	MenuItemID   string
	Price        Money
}

type PriceObservation struct {
	ID         string
	OfferID    string
	Price      Money
	ObservedAt time.Time
}

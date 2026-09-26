package entity

import (
	"encoding/json"

	money "github.com/ChristianDenniss/go-data-model/money/entity"
	userentity "github.com/ChristianDenniss/go-data-model/user/entity"
)

type Request struct {
	PlaceID            string
	FulfillmentContext FulfillmentContext
	Basket             Basket
	Filters            userentity.ComparePrefs
	Memberships        []string
	QuotePreference    string
}

type FulfillmentContext struct {
	Mode    string
	Dropoff *DropoffPoint
}

type DropoffPoint struct {
	Latitude  float64
	Longitude float64
	Label     string
}

type Basket struct {
	Lines []BasketLine
}

type BasketLine struct {
	DishID       string
	SourceItemID string
	Quantity     int
}

type Result struct {
	CompareSessionID string
	Recommendation   *PathRank
	RunnersUp        []PathRank
	AlsoConsidered   []PathRank
	UnavailablePaths []UnavailablePath
	Raw              json.RawMessage
}

type PathRank struct {
	PurchaseOptionID string
	Rank             int
	Kind             string
	Headline         string
	Confidence       string
	AllIn            money.Money
	FulfillmentMode  string
	DeliveryExecutor string
	ChannelID        string
	RationaleBullets []string
}

type UnavailablePath struct {
	PurchaseOptionID string
	Code             string
	Message          string
}

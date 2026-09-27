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
	PurchaseOptionID string       `json:"purchase_option_id"`
	Rank             int          `json:"rank"`
	Kind             string       `json:"kind"`
	Headline         string       `json:"headline"`
	Confidence       string       `json:"confidence"`
	AllIn            money.Money  `json:"all_in"`
	FulfillmentMode  string       `json:"fulfillment_mode"`
	DeliveryExecutor string       `json:"delivery_executor"`
	ChannelID        string       `json:"channel_id"`
	RationaleBullets []string     `json:"rationale_bullets"`
	ProviderID       string       `json:"provider_id,omitempty"`
	ProviderName     string       `json:"provider_name,omitempty"`
	RestaurantID     string       `json:"restaurant_id,omitempty"`
	MenuItemID       string       `json:"menu_item_id,omitempty"`
	ItemName         string       `json:"item_name,omitempty"`
	ItemSubtotal     money.Money  `json:"item_subtotal,omitempty"`
	Fees             []FeeLine    `json:"fees,omitempty"`
	DeliveryCost     money.Money  `json:"delivery_cost,omitempty"`
	ServiceFee       money.Money  `json:"service_fee,omitempty"`
	EtaMinutes       int          `json:"eta_minutes,omitempty"`
	HandoffURL       string       `json:"handoff_url,omitempty"`
}

type FeeLine struct {
	Kind   string      `json:"kind"`
	Amount money.Money `json:"amount"`
}

type UnavailablePath struct {
	PurchaseOptionID string `json:"purchase_option_id"`
	Code             string `json:"code"`
	Message          string `json:"message"`
}

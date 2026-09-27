package entity

// Hours.Service values. A restaurant's delivery hours can differ from the hours
// its doors are open (e.g. delivery stops an hour before close).
const (
	HoursStore    = "store"
	HoursDelivery = "delivery"
)

// Hours is one open interval in a restaurant's weekly schedule. DayOfWeek runs
// 0 (Sunday) to 6. Opens and Closes are "HH:MM" in the restaurant's local time;
// Closes at or before Opens means the interval runs past midnight.
type Hours struct {
	Service   string
	DayOfWeek int
	Opens     string
	Closes    string
}

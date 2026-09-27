package entity

type Account struct {
	ID             string
	Name           string
	Email          string
	Phone          string
	Role           string
	Addresses      []SavedAddress
	PaymentMethods []PaymentMethod
}

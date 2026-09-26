/**
 * Generated from go-data-model. Do not edit.
 * Regenerate: go run ./cmd/gentypes (from go-data-model).
 */

export type OrderStatus = 'pending' | 'completed' | 'cancelled'

export interface Money {
  amountCents: number
  currency: string
}

export interface Location {
  latitude: number
  longitude: number
  address: string
  city: string
  region: string
  postalCode: string
}

export interface Rating {
  average: number
  count: number
}

export interface Provider {
  id: string
  name: string
}

export interface Category {
  id: string
  slug: string
  name: string
  description: string
}

export interface Cuisine {
  id: string
  slug: string
  name: string
}

export interface Restaurant {
  id: string
  name: string
  location: Location
  cuisineIds: string[]
  categoryIds: string[]
  rating: Rating
}

export interface Item {
  id: string
  restaurantId: string
  name: string
  description: string
  section: string
}

export interface Offer {
  id: string
  restaurantId: string
  providerId: string
  menuItemId: string
  price: Money
  estimatedMinutes: number
}

export interface Observation {
  id: string
  offerId: string
  price: Money
  observedAt: string
}

export interface SavedAddress {
  id: string
  label: string
  location: Location
  current: boolean
}

export interface PaymentMethod {
  id: string
  brand: string
  last4: string
  expMonth: number
  expYear: number
  default: boolean
}

export interface Account {
  id: string
  name: string
  email: string
  phone: string
  addresses: SavedAddress[]
  paymentMethods: PaymentMethod[]
}

export interface CartLine {
  id: string
  restaurantId: string
  menuItemId: string
  providerId: string
  quantity: number
}

export interface Cart {
  id: string
  accountId: string
  lines: CartLine[]
}

export interface OrderLine {
  id: string
  menuItemId: string
  name: string
  quantity: number
}

export interface Order {
  id: string
  accountId: string
  restaurantId: string
  providerId: string
  placedAt: string
  status: OrderStatus
  total: Money
  lines: OrderLine[]
}


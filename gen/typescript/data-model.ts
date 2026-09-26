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

export interface Channel {
  id: string
  slug: string
  kind: string
  name: string
}

export interface Brand {
  id: string
  slug: string
  name: string
}

export interface Place {
  id: string
  brandId: string
  name: string
  location: Location
}

export interface Dish {
  id: string
  brandId: string
  name: string
  canonicalName: string
  description: string
}

export interface SourceStore {
  id: string
  channelId: string
  externalStoreId: string
  name: string
  location: Location
  phone: string
}

export interface SourceItem {
  id: string
  sourceCategoryId: string
  externalItemId: string
  name: string
  description: string
  available: boolean
}

export interface User {
  id: string
  name: string
  email: string
  phone: string
}

export interface ComparePrefs {
  allowedFulfillmentModes: string[]
  willingToUseAggregator: boolean
  allowedChannelIds: string[]
  blockedChannelIds: string[]
  driveThruOK: boolean
  merchantDirectOK: boolean
}

export interface CompareBasketLine {
  dishId: string
  sourceItemId: string
  quantity: number
}

export interface CompareBasket {
  lines: BasketLine[]
}

export interface CompareFulfillmentContext {
  mode: string
  dropoff: DropoffPoint | null
}

export interface CompareDropoffPoint {
  latitude: number
  longitude: number
  label: string
}

export interface CompareAPIPathRank {
  purchaseOptionId: string
  rank: number
  kind: string
  headline: string
  confidence: string
  allIn: Money
  fulfillmentMode: string
  deliveryExecutor: string
  channelId: string
  rationaleBullets: string[]
}

export interface CompareAPIUnavailablePath {
  purchaseOptionId: string
  code: string
  message: string
}

export interface CompareAPIResponse {
  compareSessionId: string
  observedAt: string
  pathsRanked: number
  recommendation: APIPathRank | null
  runnersUp: APIPathRank[]
  unavailablePaths: APIUnavailablePath[]
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


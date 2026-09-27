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

export interface Market {
  id: string
  slug: string
  name: string
  country: string
  region: string
  currency: string
  timezone: string
  status: string
  geohashPrefixes: string
}

export interface ProbeDropoff {
  id: string
  marketId: string
  label: string
  geohash: string
  location: Location
}

export interface ChannelCoverage {
  id: string
  channelId: string
  marketId: string
  status: string
  storeCount: number
  ingestRunId: string
  note: string
  lastObservedAt: string | null
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

export interface SourceMenu {
  id: string
  sourceStoreId: string
  fulfillmentMode: string
  deliveryExecutor: string
  externalMenuId: string
}

export interface SourceCategory {
  id: string
  sourceMenuId: string
  externalCategoryId: string
  name: string
  sortOrder: number
}

export interface SourceItem {
  id: string
  sourceCategoryId: string
  externalItemId: string
  name: string
  description: string
  available: boolean
  imageURL: string
}

export interface SourceMenuItemView {
  item: SourceItem
  priceCents: number
  currency: string
}

export interface SourceCategoryWithItems {
  category: SourceCategory
  items: SourceMenuItemView[]
}

export interface SourceMenuBrowse {
  store: SourceStore
  menu: SourceMenu
  categories: SourceCategoryWithItems[]
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
  lines: CompareBasketLine[]
}

export interface CompareFulfillmentContext {
  mode: string
  dropoff: CompareDropoffPoint | null
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
  recommendation: CompareAPIPathRank | null
  runnersUp: CompareAPIPathRank[]
  unavailablePaths: CompareAPIUnavailablePath[]
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

export interface RestaurantHours {
  service: string
  dayOfWeek: number
  opens: string
  closes: string
}

export interface Restaurant {
  id: string
  name: string
  location: Location
  cuisineIds: string[]
  categoryIds: string[]
  rating: Rating
  phone: string
  appURL: string
  hours: RestaurantHours[]
}

export interface Item {
  id: string
  restaurantId: string
  name: string
  description: string
  section: string
  imageURL: string
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

export interface Promotion {
  id: string
  channelId: string
  name: string
  description: string
  kind: string
  fulfillmentMode: string
  value: Money
  valueBPS: number
  startsAt: string
  endsAt: string
}

export interface PromotionConstraint {
  id: string
  promotionId: string
  minSubtotalCents: number
  code: string
  membershipRequired: boolean
  maxDiscountCents: number
}

export interface PromotionTarget {
  id: string
  promotionId: string
  placeId: string
  sourceStoreId: string
  sourceItemId: string
  dishId: string
  brandId: string
  legacyRestaurantId: string
}

export interface MembershipProduct {
  id: string
  channelId: string
  name: string
  slug: string
}

export interface Advertiser {
  id: string
  name: string
  brandId: string
  contactEmail: string
  status: string
  createdAt: string
}

export interface SponsoredCampaign {
  id: string
  advertiserId: string
  marketId: string
  name: string
  status: string
  startsAt: string
  endsAt: string
  pricingModel: string
  bidCents: number
  dailyBudgetCents: number
  totalBudgetCents: number
  currency: string
}

export interface SponsoredPlacement {
  id: string
  campaignId: string
  slot: string
  priority: number
  legacyRestaurantId: string
  placeId: string
  brandId: string
  promotionId: string
  categoryId: string
  cuisineId: string
  headline: string
  body: string
  imageURL: string
  callToAction: string
}

export interface SponsoredEvent {
  id: string
  placementId: string
  kind: string
  userId: string
  surface: string
  occurredAt: string
}

export interface SponsoredMark {
  placementId: string
  campaignId: string
  advertiserName: string
  label: string
}

export interface DealBadge {
  promotionId: string
  channelId: string
  label: string
  fulfillmentMode: string
}

export interface HomeBanner {
  id: string
  kind: string
  headline: string
  body: string
  imageURL: string
  callToAction: string
  restaurantId: string
  sponsored: SponsoredMark | null
  deal: DealBadge | null
}

export interface HomeFeedItem {
  restaurantId: string
  reason: string
  sponsored: SponsoredMark | null
  deal: DealBadge | null
}

export interface HomeSection {
  kind: string
  title: string
  items: HomeFeedItem[]
}

export interface HomeFeed {
  generatedAt: string
  banners: HomeBanner[]
  sections: HomeSection[]
}


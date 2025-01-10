export const storeGoodsStatus = {
  Selling: 'Selling',
  DemandSatisfied: 'DemandSatisfied',
  HighPrice: 'HighPrice',
  NotEnoughMinerals: 'NotEnoughMinerals',
  SpendingLimitReached: 'SpendingLimitReached',
  CapacityReached: 'CapacityReached',
  OnStrike: 'OnStrike'
} as const

export type StoreGoodsStatusEnum = (typeof storeGoodsStatus)[keyof typeof storeGoodsStatus]

export type StoreGoodsStatus = StoreGoodsStatusEnum
export const buildingStatus = {
  ConstructionStatus: 'Construction',
  ReadyStatus: 'Ready',
  ProductionStatus: 'Production',
  ResourcesNeededStatus: 'ResourcesNeeded',
  StorageNeededStatus: 'StorageNeeded'
} as const

export type BuildingStatusEnum = (typeof buildingStatus)[keyof typeof buildingStatus]

export type BuildingStatus = BuildingStatusEnum
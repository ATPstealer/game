import type { ResourceAmount } from '@/types/Resources/index.interface'

export interface ConstructBuildingPayload {
  x: number;
  y: number;
  typeId: number;
  square: number;
}

export interface Building {
  id: number;
  typeId: number;
  title: number;
  square: number;
  x: number;
  y: number;
  level: number;
  status: string;
  workStarted: string;
  workEnd: string;
  buildingGroup: string;
  buildingSubGroup: string;
  hiringNeeds: number;
  salary: number;
  workers: number;
  maxWorkers: number;
  onStrike: boolean;
}

export interface BuildingType {
  id: number;
  title: string;
  description?: string;
  cost?: number;
  requirements?: string;
  buildTime?: number;
}

export interface Blueprint {
  id: number;
  name: string;
  producedResources: ResourceAmount[];
  usedResources: ResourceAmount[];
  producedInId: number;
  productionTime: number;
}

export interface SearchBuildingParams {
  limit: -1;
  nickName?: string;
  x?: string;
  y?: string;
  buildingTypeId?: string;
}

export interface Goods {
  id: number;
  buildingId: number;
  resourceTypeId: number;
  price: number;
  sellSum: number;
  revenue: number;
  status: string;
}

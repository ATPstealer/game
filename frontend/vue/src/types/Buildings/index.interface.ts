import type { BuildingEquipment } from '@/types/Equipment/index.interface'
import type { ResourceAmount } from '@/types/Resources/index.interface'

export interface ConstructBuildingPayload {
  x: number;
  y: number;
  typeId: number;
  square: number;
}

export interface Building {
  _id: string;
  id: string;
  title: string;
  typeId: number;
  square: number;
  squareInUse: number;
  x: number;
  y: number;
  level: number;
  status?: string;
  workStarted: string;
  workEnd: string;
  hiringNeeds: number;
  salary: number;
  workers: number;
  onStrike: boolean;
  buildingType: BuildingType;
  goods: Goods[];
  buildingGroup: string;
  equipment: BuildingEquipment[];
  equipmentEffect: EquipmentEffect[];
}

export interface BuildingType {
  id: number;
  title: string;
  description: string;
  cost: number;
  requirements?: string;
  buildTime: number;
  buildingGroup: string;
  buildingSubGroup: string;
  capacity: number;
  workers: number;

}

export interface Blueprint {
  id: number;
  name: string;
  producedResources: ResourceAmount[];
  usedResources: ResourceAmount[];
  producedInId: string;
  productionTime: number;
}

export interface SearchBuildingParams {
  limit: number;
  nickName?: string;
  x?: string;
  y?: string;
  buildingTypeId?: string;
}

export interface Goods {
  resourceTypeId: number;
  price: number;
  sellSum: number;
  revenue: number;
  status?: string;
  sellStarted: string;
}

export interface EquipmentEffect {
  effectId: number;
  blueprintId: number;
  value: number;
}

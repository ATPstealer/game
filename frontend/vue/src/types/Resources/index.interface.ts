export interface Resource {
  id: number;
  resourceTypeId: number;
  amount: number;
  x: number;
  y: number;
  name: string;
  volume: number;
  weight: number;
  productionTime: number;
  producedInId: number;
}

export interface ResourceAmount {
  resourceId: number;
  amount: number;
}

export interface ResourceType {
  id: number;
  name: string;
  volume?: number;
  weight?: number;
  productionTime?: number;
  producedInId?: number;
}

export interface MarketParams {
  x?: number;
  y?: number;
  resource_type_id: number;
  trigger?: number;
}

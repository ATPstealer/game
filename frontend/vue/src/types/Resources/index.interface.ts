export interface Resource {
  _id: string;
  resourceTypeId: number;
  amount: number;
  x: number;
  y: number;
  resourceType: ResourceType;
}

export interface ResourceAmount {
  resourceId: number;
  amount: number;
}

export interface ResourceType {
  id: number;
  name: string;
  volume: number;
  weight: number;
  storeGroup: string;
}

export interface MarketParams {
  x?: number;
  y?: number;
  resourceTypeId: number;
  trigger?: number;
}

export interface ResourceMovePayload {
  resourceTypeId: number;
  amount: number;
  fromX: number;
  fromY: number;
  toX: number;
  toY: number;
}

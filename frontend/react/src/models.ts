export interface IResponse {
  data?: unknown;
  status: string;
  text: string;
}

export interface IResponseArray<T> {
  data: T[];
  status: string;
  text: string;
}

export interface IResponseEntity<T> {
  data: T;
  status: string;
  text: string;
}

export interface Characteristics {
  memory?: number;
  intelligence?: number;
  attention?: number;
  wits?: number;
  multitasking?: number;
  management?: number;
  planning?: number;
}

export interface IUser extends Characteristics{
  id?: number;
  nickName: string;
  email: string;
  password: string;
  money?: number;

}

export interface IToken {
  id?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string;
  token: string;
  ttl: number;
}

export interface IMap {
  cellName: string;
  x: number;
  y: number;
  surfaceImagePath: string;
  square: number;
  pollution: number;
  population: number;
  education: number;
  crime: number;
  medicine: number;
  elementarySchool: number;
  higherSchool: number;
}

export interface ICellOwners {
  nickName: string;
  square: number;
  x: number;
  y: number;
}

export interface ILand {
  nickName: string;
  square: number;
  x: number;
  y: number;
}

export interface IBuildings {
  id: number;
  typeId: number;
  title: number;
  square: number;
  x: number;
  y: number;
  level: number;
  status: string;
  workStarted: Date;
  workEnd: Date;
  productionId: number;
  blueprintName: string;
}

export interface IBuildingTypes {
  id: number;
  title: string;
  description: string;
  cost: number;
  requirements?: string;
  buildTime: number;
}

export interface IResources {
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

export interface IResourceTypes {
  id: number;
  name: string;
  volume: number;
  weight: number;
  productionTime: number;
  producedInId: number;
}

export interface ILogistics {
  id: number;
  resourceTypeId: number;
  amount: number;
  fromX: number;
  fromY: number;
  toX: number;
  toY: number;
  workEnd: Date;
  resourceName: string;
}

export interface IBlueprint {
  id: number;
  name: string;
  producedResources: IResourceAmount[];
  usedResources: IResourceAmount[];
  producedInId: number;
  productionTime: number;
}

export interface IResourceAmount {
  resourceId: number;
  amount: number;
}

export interface IStorage {
  id: number;
  userId: number;
  volumeOccupied: number;
  volumeMax: number;
  x: number;
  y: number;
}

export interface IOrder {
  id: number;
  userId: number;
  x: number;
  y: number;
  resourceTypeId: number;
  resourceName: string;
  amount: number;
  priceForUnit: number;
  sell: boolean;
}

export interface IGoods {
  id: number;
  buildingId: number;
  resourceTypeId: number;
  price: number;
  sellSum: number;
  revenue: number;
  sellStarted: Date;
  status: string;
}

export type Message = (message: string) => void

// TODO: почему models могут быть выше src, а конфиг нет,

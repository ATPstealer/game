export interface Characteristics {
  memory?: number;
  intelligence?: number;
  attention?: number;
  wits?: number;
  multitasking?: number;
  management?: number;
  planning?: number;
}

export interface User extends Characteristics{
  id?: number;
  nickName: string;
  email: string;
  password: string;
  money?: number;
}

export interface Land {
  nickName: string;
  square: number;
  x: number;
  y: number;
}

export interface Storage {
  id: number;
  userId: number;
  volumeOccupied: number;
  volumeMax: number;
  x: number;
  y: number;
}

export interface Order {
  id: number;
  userId: number;
  nickName: string;
  x: number;
  y: number;
  resourceTypeId: number;
  resourceName: string;
  amount: number;
  priceForUnit: number;
  sell: boolean;
  volume: number;
  weight: number;
}

export interface Logistic {
  id: number;
  resourceTypeId: number;
  amount: number;
  fromX: number;
  fromY: number;
  toX: number;
  toY: number;
  workEnd: string;
  resourceName: string;
}

export interface DataMessage {
  status: string;
  text: string;
}

export interface Params{
  limit: -1;
  nickName?: string;
  x?: string;
  y?: string;
  buildingTypeId?: string;
}

export type Coords = 'mapMaxX' | 'mapMaxY' | 'mapMinX' | 'mapMinY'

import type { ResourceType } from '@/types/Resources/index.interface'

export interface BackData {
  status: string;
  code: number;
  text?: string;
  values?: Record<string, string | number>;
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
  _id: string;
  userId: number;
  nickName: string;
  x: number;
  y: number;
  resourceTypeId: number;
  amount: number;
  priceForUnit: number;
  sell: boolean;
  resourceType: ResourceType;
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
  resourceType: ResourceType;
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

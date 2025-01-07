import type { Bank } from './Bank.ts'
import type { BuildingStatus } from './BuildingStatus.ts'
import type { BuildingType } from './BuildingType.ts'
import type { Equipment } from './Equipment.ts'
import type { EquipmentEffect } from './EquipmentEffect.ts'
import type { Goods } from './Goods.ts'
import type { Logistics } from './Logistics.ts'
import type { Production } from './Production.ts'

export interface BuildingWithData {
  /**
   * @type string
   */
  _id: string;
  /**
   * @type object | undefined
   */
  bank?: Bank;
  /**
   * @type object
   */
  buildingType: BuildingType;
  /**
   * @type array | undefined
   */
  equipment?: Equipment[];
  /**
   * @type array | undefined
   */
  equipmentEffect?: EquipmentEffect[];
  /**
   * @type array | undefined
   */
  goods?: Goods[];
  /**
   * @type integer
   */
  hiringNeeds: number;
  /**
   * @type integer
   */
  level: number;
  /**
   * @type object | undefined
   */
  logistics?: Logistics;
  /**
   * @type string
   */
  nickName: string;
  /**
   * @type boolean
   */
  onStrike: boolean;
  /**
   * @type object | undefined
   */
  production?: Production;
  /**
   * @type number
   */
  salary: number;
  /**
   * @type integer
   */
  square: number;
  /**
   * @type number
   */
  squareInUse: number;
  /**
   * @type string
   */
  status: BuildingStatus;
  /**
   * @type integer
   */
  typeId: number;
  /**
   * @type string
   */
  userId: string;
  /**
   * @type string | undefined
   */
  workEnd?: string;
  /**
   * @type string | undefined
   */
  workStarted?: string;
  /**
   * @type integer
   */
  workers: number;
  /**
   * @type integer
   */
  x: number;
  /**
   * @type integer
   */
  y: number;
}
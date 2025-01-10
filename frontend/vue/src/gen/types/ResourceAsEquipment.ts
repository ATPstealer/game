import type { EquipmentType } from './EquipmentType.ts'
import type { ResourceType } from './ResourceType.ts'

export interface ResourceAsEquipment {
  /**
   * @type string
   */
  _id: string;
  /**
   * @type number
   */
  amount: number;
  /**
   * @type object
   */
  equipmentType: EquipmentType;
  /**
   * @type object
   */
  resourceType: ResourceType;
  /**
   * @type integer
   */
  resourceTypeId: number;
  /**
   * @type string
   */
  userId: string;
  /**
   * @type integer
   */
  x: number;
  /**
   * @type integer
   */
  y: number;
}
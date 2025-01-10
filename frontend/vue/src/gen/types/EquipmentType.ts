export interface EquipmentType {
  /**
   * @type array | undefined
   */
  blueprintIds?: number[];
  /**
   * @type integer
   */
  durability: number;
  /**
   * @type integer
   */
  effectId: number;
  /**
   * @type integer
   */
  id: number;
  /**
   * @type string
   */
  name: string;
  /**
   * @type integer
   */
  resourceTypeId: number;
  /**
   * @type number
   */
  square: number;
  /**
   * @type number
   */
  value: number;
  /**
   * @type number
   */
  valueSecond: number;
}
export interface EquipmentEffect {
  /**
   * @type integer
   */
  blueprintId: number;
  /**
   * @type integer
   */
  effectId: number;
  /**
   * @type number
   */
  value: number;
  /**
   * @description Second value is considered as an average for all equipment, taking into account the first value.
   * @type number | undefined
   */
  valueSecond?: number;
}
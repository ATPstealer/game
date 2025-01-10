export interface ResourceType {
  /**
   * @type number
   */
  demand: number;
  /**
   * @type integer
   */
  id: number;
  /**
   * @type string
   */
  name: string;
  /**
   * @type string
   */
  storeGroup: string;
  /**
   * @description m3
   * @type number
   */
  volume: number;
  /**
   * @description kg
   * @type number
   */
  weight: number;
}
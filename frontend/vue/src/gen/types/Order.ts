export interface Order {
  /**
   * @type string
   */
  _id: string;
  /**
   * @type number
   */
  amount: number;
  /**
   * @type number
   */
  priceForUnit: number;
  /**
   * @type integer
   */
  resourceTypeId: number;
  /**
   * @description true - sell; false - buy
   * @type boolean
   */
  sell: boolean;
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
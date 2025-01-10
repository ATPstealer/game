export interface Storage {
  /**
   * @type string
   */
  _id: string;
  /**
   * @type string
   */
  userId: string;
  /**
   * @type number
   */
  volumeMax: number;
  /**
   * @type number
   */
  volumeOccupied: number;
  /**
   * @type integer
   */
  x: number;
  /**
   * @type integer
   */
  y: number;
}
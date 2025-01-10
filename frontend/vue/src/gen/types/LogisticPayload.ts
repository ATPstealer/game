export interface LogisticPayload {
  /**
   * @type number | undefined
   */
  amount?: number;
  /**
   * @type string | undefined
   */
  buildingId?: string;
  /**
   * @type integer | undefined
   */
  fromX?: number;
  /**
   * @type integer | undefined
   */
  fromY?: number;
  /**
   * @type integer | undefined
   */
  resourceTypeId?: number;
  /**
   * @type integer | undefined
   */
  toX?: number;
  /**
   * @type integer | undefined
   */
  toY?: number;
}
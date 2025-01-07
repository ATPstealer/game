import type { ResourceType } from './ResourceType.ts'

export interface LogisticWithData {
  /**
   * @type string | undefined
   */
  _id?: string;
  /**
   * @type number | undefined
   */
  amount?: number;
  /**
   * @type integer | undefined
   */
  fromX?: number;
  /**
   * @type integer | undefined
   */
  fromY?: number;
  /**
   * @type object | undefined
   */
  resourceType?: ResourceType;
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
  /**
   * @type string | undefined
   */
  userId?: string;
  /**
   * @type string | undefined
   */
  workEnd?: string;
}
import type { StoreGoodsStatus } from './StoreGoodsStatus.ts'

export interface Goods {
  /**
   * @type number | undefined
   */
  price?: number;
  /**
   * @type integer
   */
  resourceTypeId: number;
  /**
   * @type number | undefined
   */
  revenue?: number;
  /**
   * @type string | undefined
   */
  sellStarted?: string;
  /**
   * @type integer | undefined
   */
  sellSum?: number;
  /**
   * @type string | undefined
   */
  status?: StoreGoodsStatus;
}
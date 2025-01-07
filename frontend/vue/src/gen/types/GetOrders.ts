import type { JsonResult } from './JsonResult.ts'
import type { OrderWithData } from './OrderWithData.ts'

export interface GetOrdersQueryParams {
  /**
   * @description Order ID
   * @type string | undefined
   */
  id?: string;
  /**
   * @description User ID
   * @type string | undefined
   */
  userId?: string;
  /**
   * @description X coordinate
   * @type integer | undefined
   */
  x?: number;
  /**
   * @description Y coordinate
   * @type integer | undefined
   */
  y?: number;
  /**
   * @description Resource Type ID
   * @type integer | undefined
   */
  resourceTypeId?: number;
  /**
   * @description Sell flag
   * @type boolean | undefined
   */
  sell?: boolean;
  /**
   * @description Limit number of orders
   * @type integer | undefined
   */
  limit?: number;
  /**
   * @description Order
   * @type integer | undefined
   */
  order?: number;
  /**
   * @description Order Field
   * @type string | undefined
   */
  orderField?: string;
  /**
   * @description Page number
   * @type integer | undefined
   */
  page?: number;
}

/**
 * @description OK
 */
export type GetOrders200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: OrderWithData[];
}

/**
 * @description Internal Server Error
 */
export type GetOrders500 = JsonResult

export type GetOrdersQueryResponse = GetOrders200

export interface GetOrdersQuery {
  Response: GetOrders200;
  QueryParams: GetOrdersQueryParams;
  Errors: GetOrders500;
}
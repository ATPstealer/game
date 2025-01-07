import type { JsonResult } from './JsonResult.ts'
import type { LogisticsWithData } from './LogisticsWithData.ts'

export interface GetResourceLogisticsQueryParams {
  /**
   * @description x
   * @type integer | undefined
   */
  x?: number;
  /**
   * @description y
   * @type integer | undefined
   */
  y?: number;
  /**
   * @description Minimum capacity
   * @type number | undefined
   */
  minCapacity?: number;
}

/**
 * @description OK
 */
export type GetResourceLogistics200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: LogisticsWithData[];
}

/**
 * @description Internal Server Error
 */
export type GetResourceLogistics500 = JsonResult

export type GetResourceLogisticsQueryResponse = GetResourceLogistics200

export interface GetResourceLogisticsQuery {
  Response: GetResourceLogistics200;
  QueryParams: GetResourceLogisticsQueryParams;
  Errors: GetResourceLogistics500;
}
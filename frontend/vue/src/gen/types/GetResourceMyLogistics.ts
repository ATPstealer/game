import type { JsonResult } from './JsonResult.ts'
import type { LogisticWithData } from './LogisticWithData.ts'

/**
 * @description OK
 */
export type GetResourceMyLogistics200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: LogisticWithData[];
}

/**
 * @description Unauthorized
 */
export type GetResourceMyLogistics401 = JsonResult

/**
 * @description Internal Server Error
 */
export type GetResourceMyLogistics500 = JsonResult

export type GetResourceMyLogisticsQueryResponse = GetResourceMyLogistics200

export interface GetResourceMyLogisticsQuery {
  Response: GetResourceMyLogistics200;
  Errors: GetResourceMyLogistics401 | GetResourceMyLogistics500;
}
import type { JsonResult } from './JsonResult.ts'
import type { OrderWithData } from './OrderWithData.ts'

/**
 * @description OK
 */
export type GetMarketOrderMy200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: OrderWithData[];
}

/**
 * @description Unauthorized
 */
export type GetMarketOrderMy401 = JsonResult

/**
 * @description Internal Server Error
 */
export type GetMarketOrderMy500 = JsonResult

export type GetMarketOrderMyQueryResponse = GetMarketOrderMy200

export interface GetMarketOrderMyQuery {
  Response: GetMarketOrderMy200;
  Errors: GetMarketOrderMy401 | GetMarketOrderMy500;
}
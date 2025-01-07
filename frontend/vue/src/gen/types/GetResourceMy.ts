import type { JsonResult } from './JsonResult.ts'
import type { ResourceWithData } from './ResourceWithData.ts'

export interface GetResourceMyQueryParams {
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
}

/**
 * @description OK
 */
export type GetResourceMy200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: ResourceWithData[];
}

/**
 * @description Unauthorized
 */
export type GetResourceMy401 = JsonResult

/**
 * @description Internal Server Error
 */
export type GetResourceMy500 = JsonResult

export type GetResourceMyQueryResponse = GetResourceMy200

export interface GetResourceMyQuery {
  Response: GetResourceMy200;
  QueryParams: GetResourceMyQueryParams;
  Errors: GetResourceMy401 | GetResourceMy500;
}
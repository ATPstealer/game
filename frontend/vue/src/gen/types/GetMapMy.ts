import type { JsonResult } from './JsonResult.ts'
import type { LandLord } from './LandLord.ts'

/**
 * @description OK
 */
export type GetMapMy200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: LandLord[];
}

/**
 * @description Internal Server Error
 */
export type GetMapMy500 = JsonResult

export type GetMapMyQueryResponse = GetMapMy200

export interface GetMapMyQuery {
  Response: GetMapMy200;
  Errors: GetMapMy500;
}
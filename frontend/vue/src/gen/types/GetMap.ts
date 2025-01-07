import type { Cell } from './Cell.ts'
import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type GetMap200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: Cell[];
}

/**
 * @description Internal Server Error
 */
export type GetMap500 = JsonResult

export type GetMapQueryResponse = GetMap200

export interface GetMapQuery {
  Response: GetMap200;
  Errors: GetMap500;
}
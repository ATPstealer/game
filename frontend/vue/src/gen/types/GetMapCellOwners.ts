import type { JsonResult } from './JsonResult.ts'
import type { LandLord } from './LandLord.ts'

export interface GetMapCellOwnersQueryParams {
  /**
   * @description X coordinate
   * @type integer
   */
  x: number;
  /**
   * @description Y coordinate
   * @type integer
   */
  y: number;
}

/**
 * @description OK
 */
export type GetMapCellOwners200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: LandLord[];
}

/**
 * @description Internal Server Error
 */
export type GetMapCellOwners500 = JsonResult

export type GetMapCellOwnersQueryResponse = GetMapCellOwners200

export interface GetMapCellOwnersQuery {
  Response: GetMapCellOwners200;
  QueryParams: GetMapCellOwnersQueryParams;
  Errors: GetMapCellOwners500;
}
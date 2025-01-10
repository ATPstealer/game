import type { JsonResult } from './JsonResult.ts'
import type { LandLord } from './LandLord.ts'

/**
 * @description OK
 */
export type GetMapAllLandLords200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: LandLord[];
}

/**
 * @description Internal Server Error
 */
export type GetMapAllLandLords500 = JsonResult

export type GetMapAllLandLordsQueryResponse = GetMapAllLandLords200

export interface GetMapAllLandLordsQuery {
  Response: GetMapAllLandLords200;
  Errors: GetMapAllLandLords500;
}
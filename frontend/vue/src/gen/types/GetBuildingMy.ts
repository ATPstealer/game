import type { BuildingWithData } from './BuildingWithData.ts'
import type { JsonResult } from './JsonResult.ts'

export interface GetBuildingMyQueryParams {
  /**
   * @description Building ID to filter by
   * @type string | undefined
   */
  _id?: string;
}

/**
 * @description OK
 */
export type GetBuildingMy200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: BuildingWithData[];
}

/**
 * @description Unauthorized
 */
export type GetBuildingMy401 = JsonResult

/**
 * @description Internal Server Error
 */
export type GetBuildingMy500 = JsonResult

export type GetBuildingMyQueryResponse = GetBuildingMy200

export interface GetBuildingMyQuery {
  Response: GetBuildingMy200;
  QueryParams: GetBuildingMyQueryParams;
  Errors: GetBuildingMy401 | GetBuildingMy500;
}
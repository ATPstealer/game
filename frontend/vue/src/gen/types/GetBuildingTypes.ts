import type { BuildingType } from './BuildingType.ts'
import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type GetBuildingTypes200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: BuildingType[];
}

/**
 * @description Internal Server Error
 */
export type GetBuildingTypes500 = JsonResult

export type GetBuildingTypesQueryResponse = GetBuildingTypes200

export interface GetBuildingTypesQuery {
  Response: GetBuildingTypes200;
  Errors: GetBuildingTypes500;
}
import type { Blueprint } from './Blueprint.ts'
import type { JsonResult } from './JsonResult.ts'

export interface GetBuildingBlueprintsQueryParams {
  /**
   * @description Blueprint ID
   * @type string | undefined
   */
  id?: string;
}

/**
 * @description OK
 */
export type GetBuildingBlueprints200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: Blueprint[];
}

/**
 * @description Internal Server Error
 */
export type GetBuildingBlueprints500 = JsonResult

export type GetBuildingBlueprintsQueryResponse = GetBuildingBlueprints200

export interface GetBuildingBlueprintsQuery {
  Response: GetBuildingBlueprints200;
  QueryParams: GetBuildingBlueprintsQueryParams;
  Errors: GetBuildingBlueprints500;
}
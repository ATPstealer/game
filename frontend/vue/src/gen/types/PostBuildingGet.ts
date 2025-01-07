import type { BuildingWithData } from './BuildingWithData.ts'
import type { FindBuildingParams } from './FindBuildingParams.ts'
import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type PostBuildingGet200 = JsonResult & {
  /**
   * @type array | undefined
   */
  values?: BuildingWithData[];
}

/**
 * @description Internal Server Error
 */
export type PostBuildingGet500 = JsonResult

/**
 * @description Parameters to filter and sort buildings
 */
export type PostBuildingGetMutationRequest = FindBuildingParams

export type PostBuildingGetMutationResponse = PostBuildingGet200

export interface PostBuildingGetMutation {
  Response: PostBuildingGet200;
  Request: PostBuildingGetMutationRequest;
  Errors: PostBuildingGet500;
}
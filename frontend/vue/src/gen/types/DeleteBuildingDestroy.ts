import type { JsonResult } from './JsonResult.ts'

export interface DeleteBuildingDestroyQueryParams {
  /**
   * @description Building ID
   * @type string
   */
  _id: string;
}

/**
 * @description OK
 */
export type DeleteBuildingDestroy200 = JsonResult

/**
 * @description Unauthorized
 */
export type DeleteBuildingDestroy401 = JsonResult

/**
 * @description Internal Server Error
 */
export type DeleteBuildingDestroy500 = JsonResult

export type DeleteBuildingDestroyMutationResponse = DeleteBuildingDestroy200

export interface DeleteBuildingDestroyMutation {
  Response: DeleteBuildingDestroy200;
  QueryParams: DeleteBuildingDestroyQueryParams;
  Errors: DeleteBuildingDestroy401 | DeleteBuildingDestroy500;
}
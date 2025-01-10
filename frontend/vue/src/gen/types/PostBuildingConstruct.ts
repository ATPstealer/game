import type { ConstructBuildingPayload } from './ConstructBuildingPayload.ts'
import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type PostBuildingConstruct200 = JsonResult & {
  /**
   * @type object | undefined
   */
  values?: ConstructBuildingPayload;
}

/**
 * @description Unauthorized
 */
export type PostBuildingConstruct401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostBuildingConstruct500 = JsonResult

/**
 * @description Building construction payload
 */
export type PostBuildingConstructMutationRequest = ConstructBuildingPayload

export type PostBuildingConstructMutationResponse = PostBuildingConstruct200

export interface PostBuildingConstructMutation {
  Response: PostBuildingConstruct200;
  Request: PostBuildingConstructMutationRequest;
  Errors: PostBuildingConstruct401 | PostBuildingConstruct500;
}
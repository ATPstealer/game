import type { JsonResult } from './JsonResult.ts'
import type { StartProductionPayload } from './StartProductionPayload.ts'

/**
 * @description OK
 */
export type PostBuildingStartWork200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostBuildingStartWork401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostBuildingStartWork500 = JsonResult

/**
 * @description Start production payload
 */
export type PostBuildingStartWorkMutationRequest = StartProductionPayload

export type PostBuildingStartWorkMutationResponse = PostBuildingStartWork200

export interface PostBuildingStartWorkMutation {
  Response: PostBuildingStartWork200;
  Request: PostBuildingStartWorkMutationRequest;
  Errors: PostBuildingStartWork401 | PostBuildingStartWork500;
}
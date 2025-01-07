import type { JsonResult } from './JsonResult.ts'
import type { StartProductionPayload } from './StartProductionPayload.ts'

/**
 * @description OK
 */
export type PostBuildingStopWork200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostBuildingStopWork401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostBuildingStopWork500 = JsonResult

/**
 * @description Production stop payload
 */
export type PostBuildingStopWorkMutationRequest = StartProductionPayload

export type PostBuildingStopWorkMutationResponse = PostBuildingStopWork200

export interface PostBuildingStopWorkMutation {
  Response: PostBuildingStopWork200;
  Request: PostBuildingStopWorkMutationRequest;
  Errors: PostBuildingStopWork401 | PostBuildingStopWork500;
}
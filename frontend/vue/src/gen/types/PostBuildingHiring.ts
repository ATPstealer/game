import type { HiringPayload } from './HiringPayload.ts'
import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type PostBuildingHiring200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostBuildingHiring401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostBuildingHiring500 = JsonResult

/**
 * @description Details of hiring
 */
export type PostBuildingHiringMutationRequest = HiringPayload

export type PostBuildingHiringMutationResponse = PostBuildingHiring200

export interface PostBuildingHiringMutation {
  Response: PostBuildingHiring200;
  Request: PostBuildingHiringMutationRequest;
  Errors: PostBuildingHiring401 | PostBuildingHiring500;
}
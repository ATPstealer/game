import type { EmergencyHiringPayload } from './EmergencyHiringPayload.ts'
import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type PostBuildingEmergencyHiring200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostBuildingEmergencyHiring401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostBuildingEmergencyHiring500 = JsonResult

/**
 * @description Emergency hiring payload
 */
export type PostBuildingEmergencyHiringMutationRequest = EmergencyHiringPayload

export type PostBuildingEmergencyHiringMutationResponse = PostBuildingEmergencyHiring200

export interface PostBuildingEmergencyHiringMutation {
  Response: PostBuildingEmergencyHiring200;
  Request: PostBuildingEmergencyHiringMutationRequest;
  Errors: PostBuildingEmergencyHiring401 | PostBuildingEmergencyHiring500;
}
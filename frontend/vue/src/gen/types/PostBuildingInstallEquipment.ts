import type { InstallEquipmentPayload } from './InstallEquipmentPayload.ts'
import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type PostBuildingInstallEquipment200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostBuildingInstallEquipment401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostBuildingInstallEquipment500 = JsonResult

/**
 * @description Equipment installation payload
 */
export type PostBuildingInstallEquipmentMutationRequest = InstallEquipmentPayload

export type PostBuildingInstallEquipmentMutationResponse = PostBuildingInstallEquipment200

export interface PostBuildingInstallEquipmentMutation {
  Response: PostBuildingInstallEquipment200;
  Request: PostBuildingInstallEquipmentMutationRequest;
  Errors: PostBuildingInstallEquipment401 | PostBuildingInstallEquipment500;
}
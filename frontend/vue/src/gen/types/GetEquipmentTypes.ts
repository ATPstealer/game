import type { EquipmentType } from './EquipmentType.ts'
import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type GetEquipmentTypes200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: EquipmentType[];
}

/**
 * @description Internal Server Error
 */
export type GetEquipmentTypes500 = JsonResult

export type GetEquipmentTypesQueryResponse = GetEquipmentTypes200

export interface GetEquipmentTypesQuery {
  Response: GetEquipmentTypes200;
  Errors: GetEquipmentTypes500;
}
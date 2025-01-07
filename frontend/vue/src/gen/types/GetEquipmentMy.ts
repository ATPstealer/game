import type { JsonResult } from './JsonResult.ts'
import type { ResourceAsEquipment } from './ResourceAsEquipment.ts'

export interface GetEquipmentMyQueryParams {
  /**
   * @description x-coordinate of the equipment location
   * @type integer | undefined
   */
  x?: number;
  /**
   * @description y-coordinate of the equipment location
   * @type integer | undefined
   */
  y?: number;
}

/**
 * @description OK
 */
export type GetEquipmentMy200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: ResourceAsEquipment[];
}

/**
 * @description Unauthorized
 */
export type GetEquipmentMy401 = JsonResult

/**
 * @description Internal Server Error
 */
export type GetEquipmentMy500 = JsonResult

export type GetEquipmentMyQueryResponse = GetEquipmentMy200

export interface GetEquipmentMyQuery {
  Response: GetEquipmentMy200;
  QueryParams: GetEquipmentMyQueryParams;
  Errors: GetEquipmentMy401 | GetEquipmentMy500;
}
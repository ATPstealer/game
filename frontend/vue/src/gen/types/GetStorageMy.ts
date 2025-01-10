import type { JsonResult } from './JsonResult.ts'
import type { Storage } from './Storage.ts'

/**
 * @description OK
 */
export type GetStorageMy200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: Storage[];
}

/**
 * @description Unauthorized
 */
export type GetStorageMy401 = JsonResult

/**
 * @description Internal Server Error
 */
export type GetStorageMy500 = JsonResult

export type GetStorageMyQueryResponse = GetStorageMy200

export interface GetStorageMyQuery {
  Response: GetStorageMy200;
  Errors: GetStorageMy401 | GetStorageMy500;
}
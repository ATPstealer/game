import type { JsonResult } from './JsonResult.ts'
import type { ResourceType } from './ResourceType.ts'

/**
 * @description OK
 */
export type GetResourceTypes200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: ResourceType[];
}

/**
 * @description Internal Server Error
 */
export type GetResourceTypes500 = JsonResult

export type GetResourceTypesQueryResponse = GetResourceTypes200

export interface GetResourceTypesQuery {
  Response: GetResourceTypes200;
  Errors: GetResourceTypes500;
}
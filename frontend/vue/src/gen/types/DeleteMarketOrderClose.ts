import type { JsonResult } from './JsonResult.ts'

export interface DeleteMarketOrderCloseQueryParams {
  /**
   * @description Order ID
   * @type string
   */
  orderId: string;
}

/**
 * @description OK
 */
export type DeleteMarketOrderClose200 = JsonResult

/**
 * @description Unauthorized
 */
export type DeleteMarketOrderClose401 = JsonResult

/**
 * @description Internal Server Error
 */
export type DeleteMarketOrderClose500 = JsonResult

export type DeleteMarketOrderCloseMutationResponse = DeleteMarketOrderClose200

export interface DeleteMarketOrderCloseMutation {
  Response: DeleteMarketOrderClose200;
  QueryParams: DeleteMarketOrderCloseQueryParams;
  Errors: DeleteMarketOrderClose401 | DeleteMarketOrderClose500;
}
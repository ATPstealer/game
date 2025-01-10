import type { JsonResult } from './JsonResult.ts'
import type { LogisticsPriceParams } from './LogisticsPriceParams.ts'

/**
 * @description OK
 */
export type PostLogisticsSetPrice200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostLogisticsSetPrice401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostLogisticsSetPrice500 = JsonResult

/**
 * @description Logistics price parameters
 */
export type PostLogisticsSetPriceMutationRequest = LogisticsPriceParams

export type PostLogisticsSetPriceMutationResponse = PostLogisticsSetPrice200

export interface PostLogisticsSetPriceMutation {
  Response: PostLogisticsSetPrice200;
  Request: PostLogisticsSetPriceMutationRequest;
  Errors: PostLogisticsSetPrice401 | PostLogisticsSetPrice500;
}
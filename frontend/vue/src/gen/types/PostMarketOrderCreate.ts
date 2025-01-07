import type { JsonResult } from './JsonResult.ts'
import type { Order } from './Order.ts'

/**
 * @description OK
 */
export type PostMarketOrderCreate200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostMarketOrderCreate401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostMarketOrderCreate500 = JsonResult

/**
 * @description Order payload
 */
export type PostMarketOrderCreateMutationRequest = Order

export type PostMarketOrderCreateMutationResponse = PostMarketOrderCreate200

export interface PostMarketOrderCreateMutation {
  Response: PostMarketOrderCreate200;
  Request: PostMarketOrderCreateMutationRequest;
  Errors: PostMarketOrderCreate401 | PostMarketOrderCreate500;
}
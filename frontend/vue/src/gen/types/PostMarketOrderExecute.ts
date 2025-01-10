import type { ExecuteOrderPayload } from './ExecuteOrderPayload.ts'
import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type PostMarketOrderExecute200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostMarketOrderExecute401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostMarketOrderExecute500 = JsonResult

/**
 * @description Order execution payload
 */
export type PostMarketOrderExecuteMutationRequest = ExecuteOrderPayload

export type PostMarketOrderExecuteMutationResponse = PostMarketOrderExecute200

export interface PostMarketOrderExecuteMutation {
  Response: PostMarketOrderExecute200;
  Request: PostMarketOrderExecuteMutationRequest;
  Errors: PostMarketOrderExecute401 | PostMarketOrderExecute500;
}
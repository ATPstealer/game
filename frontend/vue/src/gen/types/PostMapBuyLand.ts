import type { BuyLandPayload } from './BuyLandPayload.ts'
import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type PostMapBuyLand200 = JsonResult & {
  /**
   * @type object | undefined
   */
  values?: BuyLandPayload;
}

/**
 * @description Internal Server Error
 */
export type PostMapBuyLand500 = JsonResult

/**
 * @description Land purchase payload
 */
export type PostMapBuyLandMutationRequest = BuyLandPayload

export type PostMapBuyLandMutationResponse = PostMapBuyLand200

export interface PostMapBuyLandMutation {
  Response: PostMapBuyLand200;
  Request: PostMapBuyLandMutationRequest;
  Errors: PostMapBuyLand500;
}
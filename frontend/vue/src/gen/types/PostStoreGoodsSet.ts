import type { JsonResult } from './JsonResult.ts'
import type { StoreGoodsPayload } from './StoreGoodsPayload.ts'

/**
 * @description OK
 */
export type PostStoreGoodsSet200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostStoreGoodsSet401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostStoreGoodsSet500 = JsonResult

/**
 * @description Store goods payload
 */
export type PostStoreGoodsSetMutationRequest = StoreGoodsPayload

export type PostStoreGoodsSetMutationResponse = PostStoreGoodsSet200

export interface PostStoreGoodsSetMutation {
  Response: PostStoreGoodsSet200;
  Request: PostStoreGoodsSetMutationRequest;
  Errors: PostStoreGoodsSet401 | PostStoreGoodsSet500;
}
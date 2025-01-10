import type { JsonResult } from './JsonResult.ts'
import type { LogisticPayload } from './LogisticPayload.ts'

/**
 * @description OK
 */
export type PostResourceMove200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostResourceMove401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostResourceMove500 = JsonResult

/**
 * @description Resource movement payload
 */
export type PostResourceMoveMutationRequest = LogisticPayload

export type PostResourceMoveMutationResponse = PostResourceMove200

export interface PostResourceMoveMutation {
  Response: PostResourceMove200;
  Request: PostResourceMoveMutationRequest;
  Errors: PostResourceMove401 | PostResourceMove500;
}
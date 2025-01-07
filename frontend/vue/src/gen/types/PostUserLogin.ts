import type { JsonResult } from './JsonResult.ts'
import type { UserPayload } from './UserPayload.ts'

/**
 * @description OK
 */
export type PostUserLogin200 = JsonResult

/**
 * @description User login payload
 */
export type PostUserLoginMutationRequest = UserPayload

export type PostUserLoginMutationResponse = PostUserLogin200

export interface PostUserLoginMutation {
  Response: PostUserLogin200;
  Request: PostUserLoginMutationRequest;
  Errors: any;
}
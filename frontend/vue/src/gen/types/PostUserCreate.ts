import type { JsonResult } from './JsonResult.ts'
import type { UserPayload } from './UserPayload.ts'

/**
 * @description OK
 */
export type PostUserCreate200 = JsonResult & {
  /**
   * @type object | undefined
   */
  values?: UserPayload;
}

/**
 * @description Internal Server Error
 */
export type PostUserCreate500 = JsonResult

/**
 * @description User creation payload
 */
export type PostUserCreateMutationRequest = UserPayload

export type PostUserCreateMutationResponse = PostUserCreate200

export interface PostUserCreateMutation {
  Response: PostUserCreate200;
  Request: PostUserCreateMutationRequest;
  Errors: PostUserCreate500;
}
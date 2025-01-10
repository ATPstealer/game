import type { JsonResult } from './JsonResult.ts'
import type { User } from './User.ts'

/**
 * @description OK
 */
export type GetUserData200 = JsonResult & {
  /**
   * @type object | undefined
   */
  data?: User;
}

/**
 * @description Unauthorized
 */
export type GetUserData401 = JsonResult

export type GetUserDataQueryResponse = GetUserData200

export interface GetUserDataQuery {
  Response: GetUserData200;
  Errors: GetUserData401;
}
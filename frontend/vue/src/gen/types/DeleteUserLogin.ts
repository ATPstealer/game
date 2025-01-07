import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type DeleteUserLogin200 = JsonResult

/**
 * @description Internal Server Error
 */
export type DeleteUserLogin500 = JsonResult

export type DeleteUserLoginMutationResponse = DeleteUserLogin200

export interface DeleteUserLoginMutation {
  Response: DeleteUserLogin200;
  Errors: DeleteUserLogin500;
}
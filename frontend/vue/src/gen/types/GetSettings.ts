import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type GetSettings200 = JsonResult

/**
 * @description Internal Server Error
 */
export type GetSettings500 = JsonResult

export type GetSettingsQueryResponse = GetSettings200

export interface GetSettingsQuery {
  Response: GetSettings200;
  Errors: GetSettings500;
}
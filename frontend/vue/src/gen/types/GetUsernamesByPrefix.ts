import type { JsonResult } from './JsonResult.ts'

export interface GetUsernamesByPrefixQueryParams {
  /**
   * @description Prefix to filter usernames
   * @type string | undefined
   */
  prefix?: string;
}

/**
 * @description OK
 */
export type GetUsernamesByPrefix200 = JsonResult

export type GetUsernamesByPrefixQueryResponse = GetUsernamesByPrefix200

export interface GetUsernamesByPrefixQuery {
  Response: GetUsernamesByPrefix200;
  QueryParams: GetUsernamesByPrefixQueryParams;
  Errors: any;
}
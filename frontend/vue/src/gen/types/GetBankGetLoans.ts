import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type GetBankGetLoans200 = JsonResult

/**
 * @description Unauthorized
 */
export type GetBankGetLoans401 = JsonResult

/**
 * @description Internal Server Error
 */
export type GetBankGetLoans500 = JsonResult

export type GetBankGetLoansQueryResponse = GetBankGetLoans200

export interface GetBankGetLoansQuery {
  Response: GetBankGetLoans200;
  Errors: GetBankGetLoans401 | GetBankGetLoans500;
}
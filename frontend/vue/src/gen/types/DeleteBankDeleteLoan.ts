import type { JsonResult } from './JsonResult.ts'

export interface DeleteBankDeleteLoanQueryParams {
  /**
   * @description Loan ID
   * @type string
   */
  _id: string;
}

/**
 * @description OK
 */
export type DeleteBankDeleteLoan200 = JsonResult

/**
 * @description Unauthorized
 */
export type DeleteBankDeleteLoan401 = JsonResult

/**
 * @description Internal Server Error
 */
export type DeleteBankDeleteLoan500 = JsonResult

export type DeleteBankDeleteLoanMutationResponse = DeleteBankDeleteLoan200

export interface DeleteBankDeleteLoanMutation {
  Response: DeleteBankDeleteLoan200;
  QueryParams: DeleteBankDeleteLoanQueryParams;
  Errors: DeleteBankDeleteLoan401 | DeleteBankDeleteLoan500;
}
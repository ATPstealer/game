import type { JsonResult } from './JsonResult.ts'
import type { RepayLoanPayload } from './RepayLoanPayload.ts'

/**
 * @description OK
 */
export type PostBankRepayLoan200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostBankRepayLoan401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostBankRepayLoan500 = JsonResult

/**
 * @description Repay loan payload
 */
export type PostBankRepayLoanMutationRequest = RepayLoanPayload

export type PostBankRepayLoanMutationResponse = PostBankRepayLoan200

export interface PostBankRepayLoanMutation {
  Response: PostBankRepayLoan200;
  Request: PostBankRepayLoanMutationRequest;
  Errors: PostBankRepayLoan401 | PostBankRepayLoan500;
}
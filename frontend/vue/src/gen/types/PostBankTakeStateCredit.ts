import type { JsonResult } from './JsonResult.ts'
import type { TakeStateCreditPayload } from './TakeStateCreditPayload.ts'

/**
 * @description OK
 */
export type PostBankTakeStateCredit200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostBankTakeStateCredit401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostBankTakeStateCredit500 = JsonResult

/**
 * @description Get state credit payload
 */
export type PostBankTakeStateCreditMutationRequest = TakeStateCreditPayload

export type PostBankTakeStateCreditMutationResponse = PostBankTakeStateCredit200

export interface PostBankTakeStateCreditMutation {
  Response: PostBankTakeStateCredit200;
  Request: PostBankTakeStateCreditMutationRequest;
  Errors: PostBankTakeStateCredit401 | PostBankTakeStateCredit500;
}
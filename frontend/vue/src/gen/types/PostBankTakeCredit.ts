import type { JsonResult } from './JsonResult.ts'
import type { TakeCreditPayload } from './TakeCreditPayload.ts'

/**
 * @description OK
 */
export type PostBankTakeCredit200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostBankTakeCredit401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostBankTakeCredit500 = JsonResult

/**
 * @description Get credit payload
 */
export type PostBankTakeCreditMutationRequest = TakeCreditPayload

export type PostBankTakeCreditMutationResponse = PostBankTakeCredit200

export interface PostBankTakeCreditMutation {
  Response: PostBankTakeCredit200;
  Request: PostBankTakeCreditMutationRequest;
  Errors: PostBankTakeCredit401 | PostBankTakeCredit500;
}
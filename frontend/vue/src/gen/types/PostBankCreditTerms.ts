import type { CreditTermsPayload } from './CreditTermsPayload.ts'
import type { JsonResult } from './JsonResult.ts'

/**
 * @description OK
 */
export type PostBankCreditTerms200 = JsonResult

/**
 * @description Unauthorized
 */
export type PostBankCreditTerms401 = JsonResult

/**
 * @description Internal Server Error
 */
export type PostBankCreditTerms500 = JsonResult

/**
 * @description Credit terms payload
 */
export type PostBankCreditTermsMutationRequest = CreditTermsPayload

export type PostBankCreditTermsMutationResponse = PostBankCreditTerms200

export interface PostBankCreditTermsMutation {
  Response: PostBankCreditTerms200;
  Request: PostBankCreditTermsMutationRequest;
  Errors: PostBankCreditTerms401 | PostBankCreditTerms500;
}
import type { CreditTermsWithData } from './CreditTermsWithData.ts'
import type { JsonResult } from './JsonResult.ts'

export interface GetBankGetCreditTermsQueryParams {
  /**
   * @description Credit limit minimum threshold
   * @type number | undefined
   */
  limit?: number;
  /**
   * @description Credit rate maximum threshold
   * @type number | undefined
   */
  rate?: number;
  /**
   * @description Credit rating maximum threshold
   * @type number | undefined
   */
  rating?: number;
}

/**
 * @description OK
 */
export type GetBankGetCreditTerms200 = JsonResult & {
  /**
   * @type array | undefined
   */
  data?: CreditTermsWithData[];
}

/**
 * @description Internal Server Error
 */
export type GetBankGetCreditTerms500 = JsonResult

export type GetBankGetCreditTermsQueryResponse = GetBankGetCreditTerms200

export interface GetBankGetCreditTermsQuery {
  Response: GetBankGetCreditTerms200;
  QueryParams: GetBankGetCreditTermsQueryParams;
  Errors: GetBankGetCreditTerms500;
}
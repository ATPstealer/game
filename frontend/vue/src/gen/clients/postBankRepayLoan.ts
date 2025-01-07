import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type {
  PostBankRepayLoanMutationRequest,
  PostBankRepayLoanMutationResponse,
  PostBankRepayLoan401,
  PostBankRepayLoan500
} from '../types/PostBankRepayLoan.ts'

/**
 * @description Pay off the loan partially or in full. Payload example {"loanId":"674ca2524dfa3a351adbf424", "Amount":122}
 * @summary Repay loan
 * {@link /bank/repay_loan}
 */
export async function postBankRepayLoan(data: PostBankRepayLoanMutationRequest, config: Partial<RequestConfig<PostBankRepayLoanMutationRequest>> = {}) {
  const res = await client<PostBankRepayLoanMutationResponse, PostBankRepayLoan401 | PostBankRepayLoan500, PostBankRepayLoanMutationRequest>({
    method: 'POST',
    url: '/bank/repay_loan',
    data,
    ...config
  })

  return res.data
}
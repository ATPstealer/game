import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostBankRepayLoanMutationRequest,
  PostBankRepayLoanMutationResponse,
  PostBankRepayLoan401,
  PostBankRepayLoan500
} from '../types/PostBankRepayLoan.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const postBankRepayLoanMutationKey = () => [{ url: '/bank/repay_loan' }] as const

export type PostBankRepayLoanMutationKey = ReturnType<typeof postBankRepayLoanMutationKey>

/**
 * @description Pay off the loan partially or in full. Payload example {"loanId":"674ca2524dfa3a351adbf424", "Amount":122}
 * @summary Repay loan
 * {@link /bank/repay_loan}
 */
async function postBankRepayLoan(data: PostBankRepayLoanMutationRequest, config: Partial<RequestConfig<PostBankRepayLoanMutationRequest>> = {}) {
  const res = await client<
    PostBankRepayLoanMutationResponse,
    ResponseErrorConfig<PostBankRepayLoan401 | PostBankRepayLoan500>,
    PostBankRepayLoanMutationRequest
  >({ method: 'POST', url: '/bank/repay_loan', data, ...config })
  
  return res.data
}

/**
 * @description Pay off the loan partially or in full. Payload example {"loanId":"674ca2524dfa3a351adbf424", "Amount":122}
 * @summary Repay loan
 * {@link /bank/repay_loan}
 */
export function usePostBankRepayLoan(
  options: {
    mutation?: MutationObserverOptions<
      PostBankRepayLoanMutationResponse,
      ResponseErrorConfig<PostBankRepayLoan401 | PostBankRepayLoan500>,
      { data: MaybeRef<PostBankRepayLoanMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBankRepayLoanMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBankRepayLoanMutationKey()

  return useMutation<
    PostBankRepayLoanMutationResponse,
    ResponseErrorConfig<PostBankRepayLoan401 | PostBankRepayLoan500>,
    { data: PostBankRepayLoanMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postBankRepayLoan(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}
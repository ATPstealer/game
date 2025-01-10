import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  DeleteBankDeleteLoanMutationResponse,
  DeleteBankDeleteLoanQueryParams,
  DeleteBankDeleteLoan401,
  DeleteBankDeleteLoan500
} from '../types/DeleteBankDeleteLoan.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const deleteBankDeleteLoanMutationKey = () => [{ url: '/bank/delete_loan' }] as const

export type DeleteBankDeleteLoanMutationKey = ReturnType<typeof deleteBankDeleteLoanMutationKey>

/**
 * @summary Delete Default Loans
 * {@link /bank/delete_loan}
 */
async function deleteBankDeleteLoan(params: DeleteBankDeleteLoanQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<DeleteBankDeleteLoanMutationResponse, DeleteBankDeleteLoan401 | DeleteBankDeleteLoan500, unknown>({
    method: 'DELETE',
    url: '/bank/delete_loan',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    params,
    ...config
  })
  
  return res.data
}

/**
 * @summary Delete Default Loans
 * {@link /bank/delete_loan}
 */
export function useDeleteBankDeleteLoan(
  options: {
    mutation?: MutationObserverOptions<
      DeleteBankDeleteLoanMutationResponse,
      DeleteBankDeleteLoan401 | DeleteBankDeleteLoan500,
      { params: MaybeRef<DeleteBankDeleteLoanQueryParams> }
    >;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? deleteBankDeleteLoanMutationKey()

  return useMutation<DeleteBankDeleteLoanMutationResponse, DeleteBankDeleteLoan401 | DeleteBankDeleteLoan500, { params: DeleteBankDeleteLoanQueryParams }>({
    mutationFn: async ({ params }) => {
      return deleteBankDeleteLoan(params, config)
    },
    mutationKey,
    ...mutationOptions
  })
}
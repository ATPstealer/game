import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type {
  DeleteBankDeleteLoanMutationResponse,
  DeleteBankDeleteLoanQueryParams,
  DeleteBankDeleteLoan401,
  DeleteBankDeleteLoan500
} from '../types/DeleteBankDeleteLoan.ts'

/**
 * @summary Delete Default Loans
 * {@link /bank/delete_loan}
 */
export async function deleteBankDeleteLoan(params: DeleteBankDeleteLoanQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<DeleteBankDeleteLoanMutationResponse, DeleteBankDeleteLoan401 | DeleteBankDeleteLoan500, unknown>({
    method: 'DELETE',
    url: '/bank/delete_loan',
    params,
    ...config
  })
  
  return res
}
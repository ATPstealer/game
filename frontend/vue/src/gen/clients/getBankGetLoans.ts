import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetBankGetLoansQueryResponse, GetBankGetLoans401, GetBankGetLoans500 } from '../types/GetBankGetLoans.ts'

/**
 * @description Return all loans connected with userId
 * @summary Get Users Loans
 * {@link /bank/get_loans}
 */
export async function getBankGetLoans(config: Partial<RequestConfig> = {}) {
  const res = await client<GetBankGetLoansQueryResponse, GetBankGetLoans401 | GetBankGetLoans500, unknown>({ method: 'GET', url: '/bank/get_loans', ...config })

  return res.data
}
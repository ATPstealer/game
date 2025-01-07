import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetBankGetCreditTermsQueryResponse, GetBankGetCreditTermsQueryParams, GetBankGetCreditTerms500 } from '../types/GetBankGetCreditTerms.ts'

/**
 * @description If defined return. Credit term where limit >= in param, rate <= in param, rating <= in param.
 * @summary Return credit terms
 * {@link /bank/get_credit_terms}
 */
export async function getBankGetCreditTerms(params?: GetBankGetCreditTermsQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetBankGetCreditTermsQueryResponse, GetBankGetCreditTerms500, unknown>({
    method: 'GET',
    url: '/bank/get_credit_terms',
    params,
    ...config
  })
  
  return res.data
}